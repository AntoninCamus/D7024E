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

func (kademlia *Kademlia) GetContacts(target *model.KademliaID, number int) []model.Contact {
	kademlia.tableMut.RLock()
	defer kademlia.tableMut.RUnlock()
	return kademlia.table.FindClosestContacts(target, number)
}

func (kademlia *Kademlia) SaveData(hash *model.KademliaID, data []byte) error {
	kademlia.filesMut.Lock()
	kademlia.files[*hash] = File{
		value:       data,
		refreshedAt: time.Now(),
		fileMut:     &sync.Mutex{},
	}
	kademlia.filesMut.Unlock()
	return nil
}

func (kademlia *Kademlia) GetData(hash *model.KademliaID) ([]byte, bool) {
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
