package kademlia

import (
	"github.com/LHJ/D7024E/kademlia/model"
	"sort"
	"sync"
	"time"
)

const parallelism = 3

// TYPES

type Kademlia struct {
	table    *model.RoutingTable
	files    map[model.KademliaID]File
	tableMut *sync.RWMutex
	filesMut *sync.RWMutex
}

type File struct {
	value       []byte
	refreshedAt time.Time
	fileMut     *sync.Mutex
}

// CONSTRUCTOR

func Init(me model.Contact) *Kademlia {
	return &Kademlia{
		table:    model.NewRoutingTable(me),
		files:    make(map[model.KademliaID]File),
		tableMut: &sync.RWMutex{},
		filesMut: &sync.RWMutex{},
	}
}

// LOCAL (THREAD SAFE, BASIC) FUNCTIONS :

func (kademlia *Kademlia) GetIdentity() model.Contact {
	return kademlia.table.Me
}

func (kademlia *Kademlia) RegisterContact(contact *model.Contact) {
	kademlia.tableMut.Lock()
	// FIXME the bucket is unlimited atm, to fix directly in it
	kademlia.table.AddContact(*contact)
	kademlia.tableMut.Unlock()
}

func (kademlia *Kademlia) getContacts(target *model.KademliaID, number int) []model.Contact {
	kademlia.tableMut.RLock()
	defer kademlia.tableMut.RUnlock()
	return kademlia.table.FindClosestContacts(target, number)
}

func (kademlia *Kademlia) saveData(data []byte, hash model.KademliaID) {
	kademlia.filesMut.Lock()
	kademlia.files[hash] = File{
		value:       data,
		refreshedAt: time.Now(),
		fileMut:     &sync.Mutex{},
	}
	kademlia.filesMut.Unlock()
}

func (kademlia *Kademlia) getData(hash *model.KademliaID) ([]byte, bool) {
	kademlia.filesMut.RLock()
	file, exists := kademlia.files[*hash]
	defer func(f File) {
		f.fileMut.Lock()
		f.refreshedAt = time.Now()
		f.fileMut.Unlock()
	}(file)
	kademlia.filesMut.RUnlock()
	return file.value, exists
}

// KADEMLIA ALGORITHMIC FUNCTIONS :

func (kademlia *Kademlia) FindContact(target *model.KademliaID) []model.Contact {
	contacts := kademlia.getContacts(target, parallelism)

	sort.Slice(contacts[:], func(i, j int) bool {
		return contacts[i].Less(&contacts[j])
	})

	//send out grpcs

	//handle the 3 incoming channels

	return contacts
}

func (kademlia *Kademlia) LookupData(target *model.KademliaID) {
	//check if present locally
	kademlia.getData(target)

	closestContacts := kademlia.getContacts(target, parallelism)
	//for i := range closestContacts {go func() {}()	} //rpc calls

	for range closestContacts { //deal with the rpc

	}

}

func (kademlia *Kademlia) StoreData(data []byte) (fileID model.KademliaID, err error) {
	// Lookup node then AddData
	targetID := model.NewKademliaID(data)
	contacts := kademlia.FindContact(targetID)

	for i := 0; i < len(contacts); i++ {
		//networking.SendStoreMessage(&(contacts[i]), data), *targetID) //this is potentially not a good idea
	}

	return *targetID, nil
}
