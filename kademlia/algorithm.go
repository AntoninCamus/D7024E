package kademlia

import (
	"errors"
	"github.com/LHJ/D7024E/kademlia/model"
	"github.com/LHJ/D7024E/kademlia/networking"
	"sort"
	"sync"
	"time"
)

const parallelism = 3
const k = 20

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
func (kademlia *Kademlia) GetIdentity() *model.Contact {
	return &kademlia.table.Me
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
  	//me := kademlia.GetIdentity()
	sort.Slice(contacts[:], func(i, j int) bool {
		return contacts[i].Less(&contacts[j])
	})


	contactedContacts := make(map[model.KademliaID]model.Contact)
	contactedContacts[*kademlia.GetIdentity().ID] = kademlia.GetIdentity()
	for _, contact := range contacts{
		contactedContacts[*contact.ID] = contact
	}

	//send out grpcs

	//handle the 3 incoming channels

	//return contacts
	return nil
}

func (kademlia *Kademlia) LookupData(fileID *model.KademliaID) ([]byte, error) {

	//check if present locally
	data, found := kademlia.GetData(fileID)
	if found {
		return data, nil
	}

	closestContacts := kademlia.GetContacts(fileID, parallelism)
	contactIn := make(chan model.Contact, parallelism)
	contactOut := make(chan model.Contact, parallelism)
	dataOut := make(chan []byte, parallelism)

	// Worker routine
	run := func(contactIn chan model.Contact, contactOut chan model.Contact, dataOut chan []byte) {
		var done = false
		for !done {
			c := <-contactIn
			if c != (model.Contact{}) {
				// Do stuff
				data, contacts, err := networking.SendFindDataMessage(&c, kademlia.GetIdentity(), fileID, 3) // Best value for nbNeighbors?
				if err != nil {
					if data != nil {
						dataOut <- data
						done = true
					} else {
						// Queue up received contacts
						for _, contact := range contacts {
							contactOut <- *contact
						}
					}
				}
			} else {
				done = true
			}
		}
	}

	// Create workers
	for i, _ := range closestContacts {
		contactIn <- closestContacts[i]
		go run(contactIn, contactOut, dataOut)
	}

	numWorkers := parallelism
	for numWorkers > 0 {
		select {
		case receivedData := <-dataOut:
			// Send special value to kill all workers
			for i := 0; i < parallelism; i++ {
				contactIn <- model.Contact{}
			}

			return receivedData, nil

		case receivedContact := <-contactOut:
			// If closer than one of closestContacts, insert it instead and add it to contactIn
			// If not insert an empty contact to kill a worker
			foundCloser := false
			for i, contact := range closestContacts {
				if receivedContact.Less(&contact) { // Check if closer
					closestContacts[i] = receivedContact

					contactIn <- receivedContact // Queue it up
					foundCloser = true
					break
				}
			}
			if !foundCloser {
				contactIn <- model.Contact{} // Kill a worker
				numWorkers--
			}
		}
	}

	return nil, errors.New("File could not be found ")
}

func (kademlia *Kademlia) StoreData(data []byte) (fileID model.KademliaID, err error) {
	/*// Lookup node then AddData
	targetID := model.NewKademliaID(data)
	closestContacts := kademlia.GetContacts(targetID, parallelism)

	for idx, c := range closestContacts { //deal with the rpc
		print(idx) //do smth with it
		print(c.Address)
	}
	*/
	targetID := model.NewKademliaID(data)

	contacts := kademlia.LookupContact(targetID)

	for _, contact := range contacts {
		networking.SendStoreMessage(&contact, kademlia.GetIdentity(), data)
	}

	return *targetID, nil
}
