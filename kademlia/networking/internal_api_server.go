package networking

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/LHJ/D7024E/kademlia/model"
)

// GrpcPort is the port where the internal API is exposed
const GrpcPort int = 9090

// InternalAPIServer is the grpc server that serves the internal API
type InternalAPIServer struct{}

// PingCall anwser to PingRequest by checking if they sent a valid KademliaID
func (s *InternalAPIServer) PingCall(ctx context.Context, in *PingRequest) (*PingAnswer, error) {
	log.Printf("Ping received")

	if len(in.GetSenderKademliaId()) != model.IDLength {
		log.Printf("Error sent : Invalid request content")
		return nil, errors.New("Invalid request content")
	}

	return &PingAnswer{ReceiverKademliaId: model.NewRandomKademliaID()[:]}, nil //FIXME return node's Kademlia ID instead
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
		grpcServer.Serve(lis)

		if err != nil {
			log.Fatal(err)
		}
	}
	go serving()
	return grpcServer
}
