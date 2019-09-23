package kademlia

import (
	"crypto/sha1"
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

func (kademlia *Kademlia) LookupData(hash string) { //find and get data

}

func (kademlia *Kademlia) Store(data []byte) {
	// Lookup node then AddData
	hash := model.NewKademliaID(data)

	targetAddr := model.NewContact(hash, "notanip") //why in the world would you need a contact
	contacts := kademlia.LookupContact(hash)

	for i := 0; i< len(contacts); i++ {
		networking.SendStoreMessage(&(contacts[i]), data, *hash) //this is potentially not a good idea
	}

}
