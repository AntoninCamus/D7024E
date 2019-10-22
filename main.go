package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

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
	ipToJoin := ""
	if len(os.Args) > 1 {
		ip := os.Args[1]
		log.Printf("Trying to join the network of %s", ip)
		ipToJoin = ip
	} else {
		log.Printf("Starting the node as first node on IP %s.", k.GetIdentity().Address)
	}

	// Start servers
	grpcSrv := kademlia.StartGrpcServer(k)
	restSrv := kademlia.StartRestServer(k, sigChan)

	// Join if needed
	if len(ipToJoin) > 0 {
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(20) // n will be between 0 and 10
		fmt.Printf("Sleeping %d seconds before joining...\n", n)
		time.Sleep(time.Duration(n) * time.Second)
		err := kademlia.JoinNetwork(k, ipToJoin)
		if err != nil {
			log.Printf("Network joining failed with error : %s", err)
			os.Exit(1)
		} else {
			log.Printf("Succeed to join the network, got contacts : %s", k.ContactStateString())
		}
	}

	// Wait for signal
	<-sigChan

	// Clean leftovers
	log.Print("Exiting servers ...")
	grpcSrv.GracefulStop()
	restSrv.Close()
}
