package kademlia

import (
	"github.com/LHJ/D7024E/kademlia/model"
	"sort"
	"sync"
	"time"
)

const parallelism = 3

// TYPES

//Kademlia is the model of the Kademlia DHT on which the algorithm works
type Kademlia struct {
	table    *model.RoutingTable
	files    map[model.KademliaID]File
	tableMut *sync.RWMutex
	filesMut *sync.RWMutex
}

//File is the internal representation of a file
type File struct {
	value       []byte
	refreshedAt time.Time
	fileMut     *sync.Mutex
}

// CONSTRUCTOR

//Init create a new kademlia object
func Init(me model.Contact) *Kademlia {
	return &Kademlia{
		table:    model.NewRoutingTable(me),
		files:    make(map[model.KademliaID]File),
		tableMut: &sync.RWMutex{},
		filesMut: &sync.RWMutex{},
	}
}

// LOCAL (THREAD SAFE, BASIC) FUNCTIONS :

//GetIdentity returns the contact information of the host
func (kademlia *Kademlia) GetIdentity() model.Contact {
	return kademlia.table.Me
}

//RegisterContact add if possible the new *contact* to the routing table
func (kademlia *Kademlia) RegisterContact(contact *model.Contact) {
	kademlia.tableMut.Lock()
	// FIXME the bucket is unlimited atm, to fix directly in it
	kademlia.table.AddContact(*contact)
	kademlia.tableMut.Unlock()
}

//GetContacts returns the *number* closest contacts to the *searchedID*
func (kademlia *Kademlia) GetContacts(searchedID *model.KademliaID, number int) []model.Contact {
	kademlia.tableMut.RLock()
	defer kademlia.tableMut.RUnlock()
	return kademlia.table.FindClosestContacts(searchedID, number)
}

//SaveData save the content of the file *content* under the *fileID*
func (kademlia *Kademlia) SaveData(fileID *model.KademliaID, content []byte) error {
	kademlia.filesMut.Lock()
	kademlia.files[*fileID] = File{
		value:       content,
		refreshedAt: time.Now(),
		fileMut:     &sync.Mutex{},
	}
	kademlia.filesMut.Unlock()
	return nil
}

//GetData returns the content corresponding to the *fileID*, as well as if the file was found
func (kademlia *Kademlia) GetData(fileID *model.KademliaID) ([]byte, bool) {
	kademlia.filesMut.RLock()
	file, exists := kademlia.files[*fileID]
	if exists {
		defer func(f File) {
			f.fileMut.Lock()
			f.refreshedAt = time.Now()
			f.fileMut.Unlock()
		}(file)
	}
	kademlia.filesMut.RUnlock()
	return file.value, exists
}

// KADEMLIA ALGORITHMIC FUNCTIONS :

func (kademlia *Kademlia) LookupContact(target *model.KademliaID) []model.Contact {
	contacts := kademlia.GetContacts(target, parallelism)

	sort.Slice(contacts[:], func(i, j int) bool {
		return contacts[i].Less(&contacts[j])
	})

	//send out grpcs

	//handle the 3 incoming channels

	return contacts
}

func (kademlia *Kademlia) LookupData(fileID *model.KademliaID) {
	//check if present locally
	kademlia.GetData(fileID)

	closestContacts := kademlia.GetContacts(fileID, parallelism)
	//for i := range closestContacts {go func() {}()	} //rpc calls

	for idx, c := range closestContacts { //deal with the rpc
		print(idx) //do smth with it
		print(c.Address)
	}

}

func (kademlia *Kademlia) StoreData(data []byte) (fileID model.KademliaID, err error) {
	// Lookup node then AddData
	targetID := model.NewKademliaID(data)
	closestContacts := kademlia.GetContacts(targetID, parallelism)

	for idx, c := range closestContacts { //deal with the rpc
		print(idx) //do smth with it
		print(c.Address)
	}

	return *targetID, nil
}
