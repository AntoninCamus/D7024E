package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/LHJ/D7024E/kademlia/model"
	"github.com/LHJ/D7024E/kademlia/networking"
)

type Kademlia struct {

	dbChan chan func() //TODO maybe change sig
	state KademliaStorage
}

type KademliaStorage struct {
	table *model.RoutingTable
	model.KademliaID
	files []Data
}

func prepare() (*Kademlia) {
	res := Kademlia{
		dbChan: nil,
		state:  KademliaStorage{},
	}
	//Run the worker

	return &res
}

type Data struct{
	hash string //how to limit length of strings :thinking:
	value string
	lastUpdatedOrAccessed int //some point in time that it was stored/retrieved, from which whether it should be automatically pruned is calculated
	keepAlive bool //default value true
}

func (kademlia *Kademlia) SaveData (data Data){
	kademlia.state.files = append(kademlia.state.files, data)
}

func (kademlia *Kademlia) LookupContact(target *model.KademliaID) []model.Contact {
	kademlia.state.table.FindClosestContacts(target, 12)
}

func (kademlia *Kademlia) LookupData(hash string) {
	//check if present locally
	newHash := model.FromString(hash)
	closestContacts := kademlia.LookupContact(newHash)
	//for i := range closestContacts {go func() {}()	} //rpc calls

	for range closestContacts{ //deal with the rpc

	}

}

func (kademlia *Kademlia) Store(data []byte) (returnHash string, err error) {
	// Lookup node then AddData
	hash := model.NewKademliaID(data)

	//targetAddr := model.NewContact(hash, "notanip") //why in the world would you need a contact
	contacts := kademlia.LookupContact(hash)

	for i := 0; i< len(contacts); i++ {
		networking.SendStoreMessage(&(contacts[i]), data)//, *hash) //this is potentially not a good idea
	}


	return hash.String(), errors.New("damn")
	//maybe the republish should be here?
}
