package kademlia

import (
	"crypto/sha1"
	"fmt"
	"github.com/LHJ/D7024E/kademlia/model"
	"log"
	"net"
	"os"
)

func getMachineHash() ([]byte, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	sha := sha1.Sum([]byte(hostname))
	return sha[:], nil
}

func getAddress() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil && ip.IsGlobalUnicast() && ip.To4() != nil {
				return ip.To4().String(), nil
			}
		}
	}
	return "", err
}

// GetContactFromHW returns a contact created from the hardware on which the node is running.
func GetContactFromHW() model.Contact {
	hash, err := getMachineHash()
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to get hostname to create contact : %s", err))
	}
	addr, err := getAddress()
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to get address to create contact : %s", err))
	}
	return model.Contact{
		ID:      model.NewKademliaID(hash),
		Address: addr,
	}
}
