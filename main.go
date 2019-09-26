package main

import (
	"github.com/LHJ/D7024E/kademlia"
	"log"
	"os"
	"os/signal"

	"github.com/LHJ/D7024E/kademlia/networking"
)

func main() {
	singleton := kademlia.Kademlia{} // Why wasn't this here?

	// Channel creation
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// Start servers
	restSrv := networking.StartRestServer(sigChan)
	grpcSrv := networking.StartGrpcServer(&singleton)

	// Wait for signal
	<-sigChan

	// Clean leftovers
	log.Print("Exiting servers ...")
	grpcSrv.GracefulStop()
	restSrv.Close()
}
