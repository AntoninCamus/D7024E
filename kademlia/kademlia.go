package kademlia

import "github.com/LHJ/D7024E/kademlia/model"

type Kademlia struct {
	Me *model.Contact
}

func (kademlia *Kademlia) RegisterContact(c *model.Contact) {
	// TODO : Add the new contact to the bucket with the algorithm of the paper
	// WARNING : do it ASYNCHRONOUSLY (return instantly) and inside the worker
}

func (kademlia *Kademlia) LookupContact(targetID *model.KademliaID, nbNeighbors int) (neighbors []model.Contact, err error) {
	// TODO
	return nil,nil
}

func (kademlia *Kademlia) LookupData(fileID *model.KademliaID) (data []byte, err error) {
	// TODO
	return nil,nil
}

func (kademlia *Kademlia) Store(data []byte) (hash KademliaID, err error) {
	// TODO
	return "", nil
}
