package kademlia

import (
	"errors"
	"github.com/LHJ/D7024E/kademlia/model"
	"log"
)

//parallelismRate is the number of parallel requests to do
const parallelismRate = 3 // Alpha

// lookupContact execute the lookupContact kademlia algorithm from the local node
func lookupContact(net *model.KademliaNetwork, target *model.KademliaID) []model.Contact {

	contactIn := make(chan model.Contact, parallelismRate)
	contactOut := make(chan model.Contact, parallelismRate)

	// Prepare the sorter of contacts
	localClosestContacts := net.GetContacts(target, parallelismRate)
	sorterClosestContacts := model.NewSorter(*target, parallelismRate)
	sorterClosestContacts.InsertContact(net.GetIdentity())
	for _, c := range localClosestContacts {
		sorterClosestContacts.InsertContact(c)
	}

	// Worker routine
	run := func(contactIn chan model.Contact, contactOut chan model.Contact) {
		var done = false
		for !done {
			c := <-contactIn //contact target
			if c != (model.Contact{}) {
				me := net.GetIdentity()
				contacts, err := sendFindContactMessage(&c, &me, target, model.BucketSize)
				//should this check whether target is you?
				if err != nil {
					log.Println("Error looking up contact")
				}

				if contacts != nil {
					for _, contact := range contacts {
						contactOut <- *contact
					}
				}

			} else {
				done = true
			}
		}
	}

	// Create workers with each of the local found closest contacts
	for _, c := range localClosestContacts {
		c.CalcDistance(target)
		contactIn <- c
		go run(contactIn, contactOut)
	}

	numWorkers := len(localClosestContacts)
	for numWorkers > 0 {
		receivedContact := <-contactOut
		foundCloser := sorterClosestContacts.InsertContact(receivedContact)
		if foundCloser {
			// If we found a closer contact, we should continue searching
			// We queue up the new found contact to the algorithm
			contactIn <- receivedContact
		} else {
			// If we did not, we should stop searching
			// We send a empty contact to kill a worker
			contactIn <- model.Contact{}
			numWorkers--
		}
	}

	return sorterClosestContacts.GetContacts()
}

// lookupData execute the LookupData kademlia algorithm from the local node
func lookupData(net *model.KademliaNetwork, fileID *model.KademliaID) ([]byte, error) {

	//check if present locally
	data, found := net.GetData(fileID)
	if found {
		return data, nil
	}

	contactIn := make(chan model.Contact, parallelismRate)
	contactOut := make(chan model.Contact, parallelismRate)
	dataOut := make(chan []byte, parallelismRate)

	// Prepare the sorter of contacts
	localClosestContacts := net.GetContacts(fileID, parallelismRate)
	sorterClosestContacts := model.NewSorter(*fileID, parallelismRate)
	sorterClosestContacts.InsertContact(net.GetIdentity())
	for _, c := range localClosestContacts {
		sorterClosestContacts.InsertContact(c)
	}

	// Worker routine
	run := func(contactIn chan model.Contact, contactOut chan model.Contact, dataOut chan []byte) {
		var done = false
		for !done {
			c := <-contactIn
			if c != (model.Contact{}) {
				// Do stuff
				me := net.GetIdentity()
				dataFound, contacts, err := sendFindDataMessage(&c, &me, fileID, 3) // Best value for nbNeighbors?
				if err != nil {
					log.Println("Error finding dataFound")
				}

				if dataFound != nil {
					dataOut <- dataFound
					done = true
				} else {
					// Queue up received contacts
					for _, contact := range contacts {
						contactOut <- *contact
					}
				}
			} else {
				done = true
			}
		}
	}

	// Create workers with each of the local found closest contacts
	for _, c := range localClosestContacts {
		c.CalcDistance(fileID)
		contactIn <- c
		go run(contactIn, contactOut, dataOut)
	}

	numWorkers := len(localClosestContacts)
	for numWorkers > 0 {
		select {
		case receivedData := <-dataOut:
			// If we found data, kill all the workers
			// Insert empty contacts to kill each worker
			for i := 0; i < numWorkers; i++ {
				contactIn <- model.Contact{}
			}
			return receivedData, nil

		case receivedContact := <-contactOut:
			// Else we receive a new contact, in this case act as LookupContact
			foundCloser := sorterClosestContacts.InsertContact(receivedContact)
			if foundCloser {
				// If we found a closer contact, we should continue searching
				// We queue up the new found contact to the algorithm
				contactIn <- receivedContact
			} else {
				// If we did not, we should stop searching
				// We send a empty contact to kill a worker
				contactIn <- model.Contact{}
				numWorkers--
			}
		}
	}

	return nil, errors.New("file could not be found ")
}

// storeData execute the StoreData kademlia algorithm from the local node
func storeData(net *model.KademliaNetwork, data []byte) (fileID model.KademliaID, err error) {
	targetID := model.NewKademliaID(data)
	contacts := lookupContact(net, targetID)

	// Store
	for _, contact := range contacts {
		me := net.GetIdentity()
		err = sendStoreMessage(&contact, &me, data)
		if err != nil {
			return fileID, err
		}
	}

	return *targetID, nil
}

// JoinNetwork execute the JoinNetwork kademlia algorithm from the local node
func JoinNetwork(net *model.KademliaNetwork, IP string) error {
	target := model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: IP,
	}

	me := net.GetIdentity()
	foundContacts, err := sendFindContactMessage(&target, &me, net.GetIdentity().ID, model.BucketSize)
	if err != nil {
		return err
	}

	for _, contact := range foundContacts {
		net.RegisterContact(contact)
	}

	return nil
}
