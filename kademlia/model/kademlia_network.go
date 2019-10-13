package model

import (
	"fmt"
	"log"
	"sync"
	"time"
)

//KademliaNetwork is the kademlia of the KademliaNetwork DHT on which the algorithm works
type KademliaNetwork struct {
	table    *RoutingTable
	files    map[KademliaID]file
	tableMut *sync.RWMutex
	filesMut *sync.RWMutex
}

//file is the internal representation of a file
type file struct {
	value       []byte
	refreshedAt time.Time
	fileMut     *sync.Mutex
}

//NewKademliaNetwork create a new kademlia object
func NewKademliaNetwork(me Contact) *KademliaNetwork {
	return &KademliaNetwork{
		table:    NewRoutingTable(me),
		files:    make(map[KademliaID]file),
		tableMut: &sync.RWMutex{},
		filesMut: &sync.RWMutex{},
	}
}

// LOCAL (THREAD SAFE, BASIC) FUNCTIONS :

//RegisterContact add if possible the new *contact* to the routing table
func (kademlia *KademliaNetwork) RegisterContact(contact *Contact) {
	if !contact.ID.equals(kademlia.GetIdentity().ID) && !kademlia.table.ContainContact(*contact.ID) {
		log.Print("Added new contact :", contact)
		kademlia.tableMut.Lock()
		kademlia.table.AddContact(*contact)
		kademlia.tableMut.Unlock()
		log.Print("Contacts known are  :", kademlia.ContactStateString())
	}
}

//GetContacts returns the *number* closest contacts to the *searchedID*
func (kademlia *KademliaNetwork) GetContacts(searchedID *KademliaID, number int) []Contact {
	kademlia.tableMut.RLock()
	defer kademlia.tableMut.RUnlock()
	return kademlia.table.FindClosestContacts(searchedID, number)
}

//SaveData save the content of the file *content* under the *fileID*
func (kademlia *KademliaNetwork) SaveData(fileID *KademliaID, content []byte) error {
	kademlia.filesMut.Lock()
	kademlia.files[*fileID] = file{
		value:       content,
		refreshedAt: time.Now(),
		fileMut:     &sync.Mutex{},
	}
	kademlia.filesMut.Unlock()
	return nil
}

//GetData returns the content corresponding to the *fileID*, as well as if the file was found
func (kademlia *KademliaNetwork) GetData(fileID *KademliaID) ([]byte, bool) {
	kademlia.filesMut.RLock()
	f, exists := kademlia.files[*fileID]
	if exists {
		defer func(f file) {
			f.fileMut.Lock()
			f.refreshedAt = time.Now()
			f.fileMut.Unlock()
		}(f)
	}
	kademlia.filesMut.RUnlock()
	return f.value, exists
}

//GetIdentity returns the contact information of the host
func (kademlia *KademliaNetwork) GetIdentity() Contact {
	return kademlia.table.GetMe()
}

func (kademlia *KademliaNetwork) ContactStateString() string {
	return kademlia.table.String()
}

func (kademlia *KademliaNetwork) FileStateString() string {
	var ret = "["
	for _, val := range kademlia.files {
		ret += fmt.Sprintf("%s,", val.value)
	}
	ret += "]"
	return ret
}
