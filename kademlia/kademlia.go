package kademlia

import "github.com/LHJ/D7024E/kademlia/model"

type Kademlia struct {
	table *RoutingTable
	KademliaID
	files []Data
}


type Data struct{
	hash string //how to limit length of strings :thinking:
	value string
	lastUpdatedOrAccessed int //some point in time that it was stored/retrieved, from which whether it should be automatically pruned is calculated
	keepAlive bool //default value true
}

func (kademlia *Kademlia) SaveData (data Data){
	kademlia.files = append(kademlia.files, data)
}

func (kademlia *Kademlia) LookupContact(target *model.Contact) {
  kademlia.table.FindClosestContacts(target.ID, 12)
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// Lookup node then AddData

}
