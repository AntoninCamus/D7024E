package networking

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

// GrpcPort is the port where the internal API is exposed
const GrpcPort int = 9090

// InternalAPIServer is the grpc server that serves the internal API
type InternalAPIServer struct {}

// StartGrpcServer start the gRPC internal API
func StartGrpcServer() *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//Creating and registering the server
	grpcServer := grpc.NewServer()

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