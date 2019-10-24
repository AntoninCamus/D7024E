package kademlia

import (
	"bytes"
	"fmt"
	"github.com/LHJ/D7024E/kademlia/model"
	"log"
	"net"
	"os"
	"sort"
)

func getHostname() (hostname string, err error) {
	hostname, err = os.Hostname()
	return
}

func getAddress() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	validIP := make([]net.IP, 0)
	for _, i := range ifaces {
		addrs, err2 := i.Addrs()
		if err2 != nil {
			return "", err2
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
				validIP = append(validIP, ip.To4())
			}
		}
	}

	if len(validIP) > 0 {
		sort.Slice(validIP, func(i, j int) bool {
			return bytes.Compare(validIP[i], validIP[j]) < 0
		})

		return validIP[0].String(),nil
	}
	
	return "", err
}

// GetContactFromHW returns a contact created from the hardware on which the node is running.
func GetContactFromHW() model.Contact {
	hostname, err := getHostname()
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to get hostname to create contact : %s", err))
	}
	addr, err := getAddress()
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to get address to create contact : %s", err))
	}
	return model.Contact{
		ID:      model.NewKademliaID([]byte(hostname)),
		Address: addr,
	}
}
