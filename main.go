package main

import (
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"

	"github.com/LHJ/D7024E/kademlia/networking"
)

type ManagementSingleton struct {
	sigChan chan os.Signal
	restSrv http.Server
	grpcSrv grpc.Server
}

func main() {
	mgmtSing := ManagementSingleton{
		sigChan: make(chan os.Signal, 1),
	}

	// TODO Servers calls should be taking as argument the signal chan and return server (or a specific structure)
	networking.StartRestServer()
	//networking.StartGrpcServer()

	// Setup signal handling
	signal.Notify(mgmtSing.sigChan, os.Interrupt)
	<-mgmtSing.sigChan // Wait for exit signal
	// Clean
}
