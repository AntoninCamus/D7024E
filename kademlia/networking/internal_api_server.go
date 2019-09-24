package networking

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/LHJ/D7024E/kademlia"
	"github.com/LHJ/D7024E/kademlia/model"

	"google.golang.org/grpc"
)

// GrpcPort is the port where the internal API is exposed
const GrpcPort int = 9090

// InternalAPIServer is the grpc server that serves the internal API
type InternalAPIServer struct {
	kademlia *kademlia.Kademlia
}

// PingCall anwser to PingRequest by checking if they sent a valid KademliaID
func (s *InternalAPIServer) PingCall(ctx context.Context, in *PingRequest) (*PingAnswer, error) {
	log.Printf("Ping received")

	if len(in.GetSenderKademliaId()) != model.IDLength {
		log.Printf("Error sent : Invalid request content")
		return nil, errors.New("Invalid request content")
	}

	return &PingAnswer{ReceiverKademliaId: s.kademlia.Me.ID[:]}, nil
}

// FindContactCall answer to FindContactRequest by sending back the NbNeighbors closest neighbors of the provided ID
func (s *InternalAPIServer) FindContactCall(ctx context.Context, in *FindContactRequest) (*FindContactAnswer, error) {
	srcContact := &model.Contact{}
	searchedID := &model.KademliaID{}

	tmpID, err := model.KademliaIDFromBytes(in.Src.ID)
	if err != nil {
		return nil, err
	}
	srcContact.ID = tmpID
	srcContact.Address = in.Src.Address
	s.kademlia.RegisterContact(srcContact)

	searchedID, err = model.KademliaIDFromBytes(in.SearchedContactId)
	if err != nil {
		return nil, err
	}
	modelContact, err := s.kademlia.LookupContact(searchedID, int(in.NbNeighbors))
	if err != nil {
		return nil, err
	}

	var newContacts []*Contact
	for _, c := range modelContact[:] {
		newContacts = append(newContacts, &Contact{
			ID:      c.ID[:],
			Address: c.Address,
		})
	}

	return &FindContactAnswer{
		Contacts: newContacts,
	}, nil
}

// FindDataCall answer to FindDataRequest by sending back the file if found, and if not act as FindContactCall
func (s *InternalAPIServer) FindDataCall(_ context.Context, in *FindDataRequest) (*FindDataAnswer, error) {
	srcContact := &model.Contact{}

	tmpID, err := model.KademliaIDFromBytes(in.Src.ID)
	if err != nil {
		return nil, err
	}
	srcContact.ID = tmpID
	srcContact.Address = in.Src.Address
	s.kademlia.RegisterContact(srcContact)

	searchedFileID, err := model.KademliaIDFromBytes(in.SearchedFileId)
	if err != nil {
		return nil, err
	}

	data, err := s.kademlia.LookupData(searchedFileID)
	if err != nil {
		//TODO Maybe we should type the error, and check for the "Not Found" particular type
		modelContacts, err := s.kademlia.LookupContact(searchedFileID, int(in.NbNeighbors))
		if err != nil {
			return nil, err
		}

		var protoContacts []*Contact
		for _, c := range modelContacts[:] {
			protoContacts = append(protoContacts, &Contact{
				ID:      c.ID[:],
				Address: c.Address,
			})
		}

		return &FindDataAnswer{
			Answer: &FindDataAnswer_DataNotFound{&FindContactAnswer{Contacts: protoContacts}},
		}, nil
	}

	return &FindDataAnswer{
		Answer: &FindDataAnswer_DataFound{data},
	}, nil
}

// StoreDataCall answer to FindDataRequest by sending back the file if found, and if not act as FindContactCall
func (s *InternalAPIServer) StoreDataCall(_ context.Context, in *StoreDataRequest) (*StoreDataAnswer, error) {
	srcContact := &model.Contact{}

	tmpID, err := model.KademliaIDFromBytes(in.Src.ID)
	if err != nil {
		return nil, err
	}
	srcContact.ID = tmpID
	srcContact.Address = in.Src.Address
	s.kademlia.RegisterContact(srcContact)

	fileID, err := s.kademlia.Store(in.Data[:])
	if err != nil {
		return &StoreDataAnswer{FileId: make([]byte, 0)}, nil
	}

	return &StoreDataAnswer{FileId: fileID}, nil
}

// StartGrpcServer start the gRPC internal API
func StartGrpcServer() *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//Creating and registering the server
	grpcServer := grpc.NewServer()
	RegisterInternalApiServiceServer(grpcServer, &InternalAPIServer{})

	log.Printf("GrpcServer ready ...")
	serving := func() {
		//Blocking call
		err = grpcServer.Serve(lis)

		if err != nil {
			log.Fatal(err)
		}
	}
	go serving()
	return grpcServer
}
