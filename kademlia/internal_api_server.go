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
		log.Printf("GRPC SERV : Error sent : Invalid request content")
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

func parseAndRegisterContact(k *model.KademliaNetwork, unparsedContact *pb.Contact) (*model.Contact, error) {
	parsedContact := &model.Contact{}

	tmpID, err := model.KademliaIDFromBytes(unparsedContact.ID)
	if err != nil {
		return nil, err
	}
	parsedContact.ID = tmpID
	parsedContact.Address = unparsedContact.Address
	k.RegisterContact(parsedContact)
	return parsedContact, err
}

// getDataAndConvert returns the *nb* closest contact to *ID* present in *network*
func getDataAndConvert(network *model.KademliaNetwork, searchedID *model.KademliaID, nb int) []*pb.Contact {
	modelContact := network.GetContacts(searchedID, nb)
	// Before returning any contact we check their availability
	allCheckPing(network, modelContact)

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
	_, err := parseAndRegisterContact(s.kademlia, in.Src)
	if err != nil {
		return nil, err
	}

	searchedID, err := model.KademliaIDFromBytes(in.SearchedContactId)
	if err != nil {
		return nil, err
	}

	contacts := getDataAndConvert(s.kademlia, searchedID, int(in.NbNeighbors))

	return &pb.FindContactAnswer{
		Contacts: contacts,
	}, nil
}

// FindDataCall answer to FindDataRequest by sending back the file if found, and if not act as FindContactCall
func (s *internalAPIServer) FindDataCall(_ context.Context, in *pb.FindDataRequest) (*pb.FindDataAnswer, error) {
	srcContact, err := parseAndRegisterContact(s.kademlia, in.Src)
	if err != nil {
		return nil, err
	}

	searchedFileID, err := model.KademliaIDFromBytes(in.SearchedFileId)
	if err != nil {
		return nil, err
	}

	data, found := s.kademlia.GetData(searchedFileID)
	if !found {
		contacts := getDataAndConvert(s.kademlia, searchedFileID, int(in.NbNeighbors))

		log.Printf("GRPC SERV : Data not found, sending back %d contacts to %s", len(contacts), srcContact.String())
		return &pb.FindDataAnswer{
			Answer: &pb.FindDataAnswer_DataNotFound{&pb.FindContactAnswer{Contacts: contacts}},
		}, nil
	}

	log.Printf("GRPC SERV : Sending back data %s to %s", data, srcContact.String())
	return &pb.FindDataAnswer{
		Answer: &pb.FindDataAnswer_DataFound{data},
	}, nil
}

// StoreDataCall store the file, and send if the store was done
func (s *internalAPIServer) StoreDataCall(_ context.Context, in *pb.StoreDataRequest) (*pb.StoreDataAnswer, error) {
	_, err := parseAndRegisterContact(s.kademlia, in.Src)
	if err != nil {
		return nil, err
	}

	fileID := model.NewKademliaID(in.Data)
	err = s.kademlia.SaveData(fileID, in.Data[:])
	if err != nil {
		return &pb.StoreDataAnswer{Ok: false}, nil
	}
	log.Printf("GRPC SERV : Data '%s' stored", in.Data)
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
