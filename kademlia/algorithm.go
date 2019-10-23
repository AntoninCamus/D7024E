package kademlia

import (
	"errors"
	"github.com/LHJ/D7024E/kademlia/model"
	"log"
	"sort"
)

//parallelismRate is the number of parallel requests to do
const parallelismRate = 3 // Alpha

// lookupContact execute the lookupContact kademlia algorithm from the local node
func lookupContact(net *model.KademliaNetwork, target *model.KademliaID) []model.Contact {

	contactIn := make(chan model.Contact, parallelismRate*model.BucketSize)
	contactOut := make(chan []*model.Contact, parallelismRate*model.BucketSize)

	// Prepare the sorter of contacts
	localClosestContacts := net.GetContacts(target, parallelismRate)
	sorterClosestContacts := model.NewSorter(*target, parallelismRate)
	sorterClosestContacts.InsertContact(net.GetIdentity())
	for _, c := range localClosestContacts {
		sorterClosestContacts.InsertContact(c)
	}

	// Worker routine
	run := func(contactIn chan model.Contact, contactOut chan []*model.Contact) {
		var done = false
		for !done {
			c := <-contactIn //contact target
			log.Printf("Worker received %s\n", c.String())
			if c != (model.Contact{}) {
				me := net.GetIdentity()
				contacts, err := sendFindContactMessage(&c, &me, target, model.BucketSize)
				//should this check whether target is you?
				if err != nil {
					log.Println("Error looking up contact")
				}

				if contacts != nil {
					contactOut <- contacts
					//log.Printf("Worker sent %s\n", c.String())
				}

			} else {
				done = true
			}
		}
	}

	numWorkers := 0

	// Create workers with each of the local found closest contacts
	for _, c := range localClosestContacts {
		c.CalcDistance(target)
		contactIn <- c
		go run(contactIn, contactOut)
		numWorkers++
	}

	for numWorkers > 0 {
		receivedContacts := <-contactOut

		// Sort the contacts for insertion
		for _, c := range receivedContacts {
			c.CalcDistance(target)
		}
		sort.Slice(receivedContacts, func(i, j int) bool {
			return receivedContacts[j].Less(receivedContacts[i])
		})

		// If we found a closer contact, we should continue searching
		// We queue up the new found contact to the algorithm
		foundCloser := false
		for _, c := range receivedContacts {
			isCloser := sorterClosestContacts.InsertContact(*c)
			if isCloser {
				foundCloser = true
				contactIn <- *c
			} else {
				break
			}
		}
		if !foundCloser {
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

	contactIn := make(chan model.Contact, parallelismRate*model.BucketSize)
	contactOut := make(chan []*model.Contact, parallelismRate*model.BucketSize)
	dataOut := make(chan []byte, parallelismRate)

	// Prepare the sorter of contacts
	localClosestContacts := net.GetContacts(fileID, parallelismRate)
	sorterClosestContacts := model.NewSorter(*fileID, parallelismRate)
	sorterClosestContacts.InsertContact(net.GetIdentity())
	for _, c := range localClosestContacts {
		sorterClosestContacts.InsertContact(c)
	}

	// Worker routine
	run := func(contactIn chan model.Contact, contactOut chan []*model.Contact, dataOut chan []byte) {
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
					contactOut <- contacts
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

		case receivedContacts := <-contactOut:
			// Sort the contacts for insertion
			for _, c := range receivedContacts {
				c.CalcDistance(fileID)
			}
			sort.Slice(receivedContacts, func(i, j int) bool {
				return receivedContacts[j].Less(receivedContacts[i])
			})

			// If we found a closer contact, we should continue searching
			// We queue up the new found contact to the algorithm
			foundCloser := false
			for _, c := range receivedContacts {
				isCloser := sorterClosestContacts.InsertContact(*c)
				if isCloser {
					foundCloser = true
					contactIn <- *c
				} else {
					break
				}
			}
			if !foundCloser {
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
func storeData(net *model.KademliaNetwork, data []byte) (model.KademliaID, error) {
	me := net.GetIdentity()
	success := make(chan bool, parallelismRate)
	targetID := model.NewKademliaID(data)
	contacts := lookupContact(net, targetID)

	// Store func call
	store := func(dst model.Contact, src model.Contact, data []byte, success chan bool) {
		err := sendStoreMessage(&dst, &src, data)
		if err != nil {
			log.Printf("error, store of %b on %s failed : %s", data, dst.String(), err.Error())
		}
		success <- err == nil
	}

	nbWorkers := 0
	for _, contact := range contacts {
		if contact.ID != nil {
			nbWorkers++
			store(contact, me, data, success)
		}
	}
	for i := nbWorkers; i > 0; i-- {
		res := <-success
		if res {
			return *targetID, nil
		}
	}

	return *targetID, errors.New("could not store the file on any node")
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
