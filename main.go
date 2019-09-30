package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/LHJ/D7024E/kademlia"
	"github.com/LHJ/D7024E/kademlia/model"
)

func main() {
	// Channel creation
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	k := model.NewKademliaNetwork(kademlia.GetContactFromHW())

	// Start servers
	restSrv := kademlia.StartRestServer(k, sigChan)
	grpcSrv := kademlia.StartGrpcServer(k)
	// Wait for signal
	<-sigChan

	// Clean leftovers
	log.Print("Exiting servers ...")
	grpcSrv.GracefulStop()
	restSrv.Close()
}
