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
	if len(in.GetSenderKademliaId()) != model.IDLength {
		log.Printf("GRPC : Error sent : Invalid request content")
		return nil, errors.New("invalid request content")
	}

	return &pb.PingAnswer{ReceiverKademliaId: s.kademlia.GetIdentity().ID[:]}, nil
}

// allCheckPing ping every contact and removes from the network distant one
func allCheckPing(k *model.KademliaNetwork, contacts []model.Contact) bool {
	result := true
	for _, c := range contacts[:] {
		if !sendPingMessage(&c, false) {
			k.UnregisterContact(c)
			result = false
		}
	}
	return result
}

// getDataAndConvert returns the *nb* closest contact to *ID* present in *network*
func getDataAndConvert(network *model.KademliaNetwork, ID *model.KademliaID, nb int) []*pb.Contact {
	modelContact := network.GetContacts(ID, nb)
	// Before returning any contact we check their availability
	allCheckPing(network, modelContact)
	/*for allCheckPing(network, modelContact) {
		// WARN : will create a pingstorm, and that is not that bad to not have exactly 20 contacts
		// If allCheck removed a contact, we retry to return the correct number
		modelContact = network.GetContacts(ID, nb)
	}*/

	// Then we convert them into *pb.Contact
	var pbContacts []*pb.Contact
	for _, c := range modelContact[:] {
		pbContacts = append(pbContacts, &pb.Contact{
			ID:      c.ID[:],
			Address: c.Address,
		})
	}

	return pbContacts
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
	s.kademlia.RegisterContact(srcContact)

	contacts := getDataAndConvert(s.kademlia, searchedID, int(in.NbNeighbors))

	return &pb.FindContactAnswer{
		Contacts: contacts,
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
		contacts := getDataAndConvert(s.kademlia, searchedFileID, int(in.NbNeighbors))

		log.Printf(" GRPC : Data not found, sending back %d contacts to %s", len(contacts), srcContact.String())
		return &pb.FindDataAnswer{
			Answer: &pb.FindDataAnswer_DataNotFound{&pb.FindContactAnswer{Contacts: contacts}},
		}, nil
	}

	log.Printf("GRPC : Sending back data %s to %s", data, srcContact.String())
	return &pb.FindDataAnswer{
		Answer: &pb.FindDataAnswer_DataFound{data},
	}, nil
}

// StoreDataCall store the file, and send if the store was done
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
	log.Printf("GRPC : Storing data %s", srcContact.String())
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
