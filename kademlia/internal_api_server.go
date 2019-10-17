package kademlia

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/LHJ/D7024E/kademlia/model"
	pb "github.com/LHJ/D7024E/protogen"
	"google.golang.org/grpc"
)

// GrpcPort is the port where the internal API is exposed
const grpcPort int = 9090

// InternalAPIServer is the grpc server that serves the internal API
type internalAPIServer struct {
	kademlia *model.KademliaNetwork
}

// PingCall anwser to PingRequest by checking if they sent a valid KademliaID
func (s *internalAPIServer) PingCall(ctx context.Context, in *pb.PingRequest) (*pb.PingAnswer, error) {
	log.Printf("Ping received")

	if len(in.GetSenderKademliaId()) != model.IDLength {
		log.Printf("Error sent : Invalid request content")
		return nil, errors.New("invalid request content")
	}

	return &pb.PingAnswer{ReceiverKademliaId: s.kademlia.GetIdentity().ID[:]}, nil
}

// FindContactCall answer to FindContactRequest by sending back the NbNeighbors closest neighbors of the provided ID
func (s *internalAPIServer) FindContactCall(ctx context.Context, in *pb.FindContactRequest) (*pb.FindContactAnswer, error) {
	srcContact := &model.Contact{}
	searchedID := &model.KademliaID{}

	tmpID, err := model.KademliaIDFromBytes(in.Src.ID)
	if err != nil {
		return nil, err
	}
	srcContact.ID = tmpID
	srcContact.Address = in.Src.Address

	searchedID, err = model.KademliaIDFromBytes(in.SearchedContactId)
	if err != nil {
		return nil, err
	}
	modelContact := s.kademlia.GetContacts(searchedID, int(in.NbNeighbors))

	s.kademlia.RegisterContact(srcContact)

	var newContacts []*pb.Contact
	for _, c := range modelContact[:] {
		newContacts = append(newContacts, &pb.Contact{
			ID:      c.ID[:],
			Address: c.Address,
		})
	}

	return &pb.FindContactAnswer{
		Contacts: newContacts,
	}, nil
}

// FindDataCall answer to FindDataRequest by sending back the file if found, and if not act as FindContactCall
func (s *internalAPIServer) FindDataCall(_ context.Context, in *pb.FindDataRequest) (*pb.FindDataAnswer, error) {
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

		var protoContacts []*pb.Contact
		for _, c := range modelContacts[:] {
			protoContacts = append(protoContacts, &pb.Contact{
				ID:      c.ID[:],
				Address: c.Address,
			})
		}

		return &pb.FindDataAnswer{
			Answer: &pb.FindDataAnswer_DataNotFound{&pb.FindContactAnswer{Contacts: protoContacts}},
		}, nil
	}

	return &pb.FindDataAnswer{
		Answer: &pb.FindDataAnswer_DataFound{data},
	}, nil
}

// StoreDataCall answer to FindDataRequest by sending back the file if found, and if not act as FindContactCall
func (s *internalAPIServer) StoreDataCall(_ context.Context, in *pb.StoreDataRequest) (*pb.StoreDataAnswer, error) {
	srcContact := &model.Contact{}

	tmpID, err := model.KademliaIDFromBytes(in.Src.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid sender ID : %s", err)
	}
	srcContact.ID = tmpID
	srcContact.Address = in.Src.Address
	s.kademlia.RegisterContact(srcContact)

	fileID := model.NewKademliaID(in.Data)
	err = s.kademlia.SaveData(fileID, in.Data[:])
	if err != nil {
		return &pb.StoreDataAnswer{Ok: false}, nil
	}
	return &pb.StoreDataAnswer{Ok: true}, nil
}

// StartGrpcServer start the gRPC internal API
func StartGrpcServer(kademlia *model.KademliaNetwork) *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("gRPC server could not start, failed to listen : %v", err)
	}

	//Creating and registering the server
	grpcServer := grpc.NewServer()
	pb.RegisterInternalApiServiceServer(grpcServer, &internalAPIServer{kademlia: kademlia})

	serving := func() {
		//Blocking call
		err = grpcServer.Serve(lis)

		if err != nil {
			log.Fatal(err)
		}
	}
	go serving()
	log.Printf("gRPC server is ready")

	return grpcServer
}
