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
		log.Printf("Error sent : Invalid request content")
		return nil, errors.New("invalid request content")
	}

	return &pb.PingAnswer{ReceiverKademliaId: s.kademlia.GetIdentity().ID[:]}, nil
}

// parallelAllCheckPing ping every contact provided and removes from the network distant one
/*func parallelAllCheckPing(k *model.KademliaNetwork, contacts []model.Contact) bool {
	// Prepare channels
	contactIn := make(chan model.Contact, len(contacts))
	resultOut := make(chan bool, len(contacts))
	for _,c := range contacts {
		contactIn <- c
	}

	// Worker routine
	pinger := func(kl *model.KademliaNetwork, contactIn chan model.Contact, resultOut chan bool) {
		var done = false
		for !done {
			c := <-contactIn
			log.Printf("Ping worker received %s\n", c.String())
			if c != (model.Contact{}) {
				result := sendPingMessage(&c)
				if ! result {
					kl.UnregisterContact(c)
				}
				resultOut <- result
			} else {
				done = true
			}
		}
	}

	for i := 0; i < len(contacts) && i < parallelismRate; i++ {
		go pinger(k, contactIn, resultOut)
	}

	result := true
	for _ = range contacts {
		res := <- resultOut
		if ! res {
			result = false
		}
	}

	for i := 0; i < len(contacts) && i < parallelismRate; i++ {
		contactIn <- model.Contact{}
	}

	return result
}
*/

// allCheckPing ping every contact and removes from the network distant one
func allCheckPing(k *model.KademliaNetwork, contacts []model.Contact) bool {
	result := true
	for _, c := range contacts[:] {
		if !sendPingMessage(&c) {
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
		// WARN : will create a pingstorm, maybe it is preferable to return less contacts
		// If allCheck removed a contact, we retry to return the correct number
		modelContact = network.GetContacts(ID, nb)
	}*/

	// Then we convert them into *pb.Contact
	var pbContacts []*pb.Contact
	for _, c := range modelContact[:] {
		if sendPingMessage(&c) {
			pbContacts = append(pbContacts, &pb.Contact{
				ID:      c.ID[:],
				Address: c.Address,
			})
		} else {
			network.UnregisterContact(c)
		}
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

		return &pb.FindDataAnswer{
			Answer: &pb.FindDataAnswer_DataNotFound{&pb.FindContactAnswer{Contacts: contacts}},
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
