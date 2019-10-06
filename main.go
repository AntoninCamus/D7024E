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

	// Parse arguments
	if len(os.Args) > 1 {
		ip := os.Args[1]
		log.Printf("Trying to join the network of %s", ip)
		err := kademlia.JoinNetwork(k, ip)
		if err != nil {
			log.Printf("Failed with error : %s", err)
			os.Exit(1)
		} else {
			log.Printf("Succeed to join the network, got contacts : %s", k.PrintContactState())
		}
	} else {
		log.Println("Starting the node as first node.")
	}

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
