package main

import (
	"github.com/LHJ/D7024E/kademlia"
	"log"
	"os"
	"os/signal"

	"github.com/LHJ/D7024E/kademlia/networking"
)

func main() {
	// Channel creation
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	k := kademlia.Init(kademlia.GetContactFromHW())

	// Start servers
	restSrv := networking.StartRestServer(sigChan)
	grpcSrv := networking.StartGrpcServer(k)
	// Wait for signal
	<-sigChan

	// Clean leftovers
	log.Print("Exiting servers ...")
	grpcSrv.GracefulStop()
	restSrv.Close()
}
