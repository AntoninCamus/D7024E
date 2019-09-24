package kademlia

import (
	"errors"
	"github.com/LHJ/D7024E/kademlia/model"
	"github.com/LHJ/D7024E/kademlia/networking"
	"time"
)

const alpha = 3

type Kademlia struct {
	dbChan chan func() //TODO maybe change sig
	table  *model.RoutingTable
	files  map[model.KademliaID][]byte
}

func Init(me model.Contact) *Kademlia {
	dbChan := make(chan func(), 100)

	//Run the worker
	go func(c chan func()) {
		for true {
			f := <-c
			f()
		}
	}(dbChan)

	return &Kademlia{
		dbChan: dbChan,
		table:  model.NewRoutingTable(me),
		files:  make(map[model.KademliaID][]byte),
	}
}

type Data struct {
	hash                  string //how to limit length of strings :thinking:
	value                 string
	lastUpdatedOrAccessed time.Time //some point in time that it was stored/retrieved, from which whether it should be automatically pruned is calculated
	keepAlive             bool      //default value true
}

func (kademlia *Kademlia) SaveData(data []byte, hash model.KademliaID) {
	kademlia.files[hash] = data
}

func (kademlia *Kademlia) GetContact(target *model.KademliaID) []model.Contact {
	// FIXME : Make thread safe
	return kademlia.table.FindClosestContacts(target, alpha)
}

func (kademlia *Kademlia) FindContact(target *model.KademliaID) []model.Contact {
	contacts := kademlia.table.FindClosestContacts(target, alpha)

	var contactsArray []model.Contact
	for i, contact := range contacts{ // sort that shit
		contact.CalcDistance(kademlia.table.Me.ID)
		contactsArray[i] = contact
	}


	return contacts
}

func (kademlia *Kademlia) GetData(hash model.KademliaID) ([]byte, bool) {
	// FIXME : Make thread safe
	data, exists := kademlia.files[hash]
	if !exists { //if file doesn't exist
		return nil, exists
	}
	return data, exists
}

	func (kademlia *Kademlia) LookupData(hash model.KademliaID) {
	//check if present locally
	kademlia.GetData(hash)


	closestContacts := kademlia.getContact(&hash)
	//for i := range closestContacts {go func() {}()	} //rpc calls

	for range closestContacts { //deal with the rpc

	}

}

func (kademlia *Kademlia) Store(data []byte) (returnHash string, err error) {
	// Lookup node then AddData
	hash := model.NewKademliaID(data)

	//targetAddr := model.NewContact(hash, "notanip") //why in the world would you need a contact
	contacts := kademlia.LookupContact(hash)

	for i := 0; i < len(contacts); i++ {
		networking.SendStoreMessage(&(contacts[i]), data) //, *hash) //this is potentially not a good idea
	}

	return hash.String(), errors.New("damn")
	//maybe the republish should be here?
}
