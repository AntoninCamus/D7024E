package model

import (
	"fmt"
	"log"
	"sync"
	"time"
)

//KademliaNetwork is the kademlia of the KademliaNetwork DHT on which the algorithm works
type KademliaNetwork struct {
	table    *routingTable
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
		table:    newRoutingTable(me),
		files:    make(map[KademliaID]file),
		tableMut: &sync.RWMutex{},
		filesMut: &sync.RWMutex{},
	}
}

// LOCAL (THREAD SAFE, BASIC) FUNCTIONS :

//RegisterContact add if possible the new *contact* to the routing table
func (kademlia *KademliaNetwork) RegisterContact(contact *Contact) {
	if !contact.ID.equals(kademlia.GetIdentity().ID) {
		kademlia.tableMut.Lock()
		kademlia.table.addContact(*contact)
		kademlia.tableMut.Unlock()
		if !!kademlia.table.containContact(*contact.ID) {
			log.Printf("Added new contact %s,\n new state is %s", contact.String(), kademlia.ContactStateString())
		}
	}
}

//UnregisterContact remove the *contact* from the routing table
func (kademlia *KademliaNetwork) UnregisterContact(contact Contact) {
	kademlia.tableMut.Lock()
	kademlia.table.removeContact(contact)
	kademlia.tableMut.Unlock()
	log.Printf("Removed contact %s,\n new state is %s", contact.String(), kademlia.ContactStateString())
}

//GetContacts returns the *number* closest contacts to the *searchedID*
func (kademlia *KademliaNetwork) GetContacts(searchedID *KademliaID, number int) []Contact {
	kademlia.tableMut.RLock()
	defer kademlia.tableMut.RUnlock()
	return kademlia.table.findClosestContacts(searchedID, number)
}

//SaveData save the content of the file *content* under the *fileID*
func (kademlia *KademliaNetwork) SaveData(fileID *KademliaID, content []byte) error {
	kademlia.filesMut.RLock()
	if kademlia.files[*fileID].value == nil {
		log.Printf("file %s added, new state is %s", fileID, kademlia.fileStateString())
	} else {
		log.Printf("file %s updated, new state is %s", fileID, kademlia.fileStateString())
	}
	kademlia.filesMut.RUnlock()
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
	return kademlia.table.getMe()
}

//ContactStateString return the routing table internal state on the form of a string
func (kademlia *KademliaNetwork) ContactStateString() string {
	return kademlia.table.String()
}

//fileStateString return the files table state on the form of a string
func (kademlia *KademliaNetwork) fileStateString() string {
	var ret = "["
	for _, val := range kademlia.files {
		ret += fmt.Sprintf("%s,", val.value)
	}
	ret += "]"
	return ret
}
