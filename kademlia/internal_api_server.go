package kademlia

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/LHJ/D7024E/kademlia/model"
	"google.golang.org/grpc"
)

// GrpcPort is the port where the internal API is exposed
const grpcPort int = 9090

// InternalAPIServer is the grpc server that serves the internal API
type internalAPIServer struct {
	kademlia *model.KademliaNetwork
}

// PingCall anwser to PingRequest by checking if they sent a valid KademliaID
func (s *internalAPIServer) PingCall(ctx context.Context, in *PingRequest) (*PingAnswer, error) {
	log.Printf("Ping received")

	if len(in.GetSenderKademliaId()) != model.IDLength {
		log.Printf("Error sent : Invalid request content")
		return nil, errors.New("invalid request content")
	}

	return &PingAnswer{ReceiverKademliaId: s.kademlia.GetIdentity().ID[:]}, nil
}

// FindContactCall answer to FindContactRequest by sending back the NbNeighbors closest neighbors of the provided ID
func (s *internalAPIServer) FindContactCall(ctx context.Context, in *FindContactRequest) (*FindContactAnswer, error) {
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
	modelContact := s.kademlia.GetContacts(searchedID, int(in.NbNeighbors))

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
func (s *internalAPIServer) FindDataCall(_ context.Context, in *FindDataRequest) (*FindDataAnswer, error) {
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

	data, found := s.kademlia.GetData(searchedFileID)
	if !found {
		modelContacts := s.kademlia.GetContacts(searchedFileID, int(in.NbNeighbors))

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
func (s *internalAPIServer) StoreDataCall(_ context.Context, in *StoreDataRequest) (*StoreDataAnswer, error) {
	srcContact := &model.Contact{}

	tmpID, err := model.KademliaIDFromBytes(in.Src.ID)
	if err != nil {
		return nil, fmt.Errorf("Invalid sender ID : %s", err)
	}
	srcContact.ID = tmpID
	srcContact.Address = in.Src.Address
	s.kademlia.RegisterContact(srcContact)

	fileID := model.NewKademliaID(in.Data)
	err = s.kademlia.SaveData(fileID, in.Data[:])
	if err != nil {
		return &StoreDataAnswer{Ok: false}, nil
	}
	return &StoreDataAnswer{Ok: true}, nil
}

// StartGrpcServer start the gRPC internal API
func StartGrpcServer(kademlia *model.KademliaNetwork) *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//Creating and registering the server
	grpcServer := grpc.NewServer()
	RegisterInternalApiServiceServer(grpcServer, &internalAPIServer{kademlia: kademlia})

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
