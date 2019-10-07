package kademlia

import (
	"errors"
	"github.com/LHJ/D7024E/kademlia/model"
)

const parallelism = 3
const k = 20

func LookupContact(net *model.KademliaNetwork, target *model.KademliaID) []model.Contact {
	//check if present locally
	//if target == net.GetIdentity().ID {
	//	return nil
	//}

	closestContacts := net.GetContacts(target, parallelism)
	contactIn := make(chan model.Contact, parallelism)
	contactOut := make(chan model.Contact, parallelism)

	// Worker routine
	run := func(contactIn chan model.Contact, contactOut chan model.Contact) {
		var done = false
		for !done {
			c := <-contactIn //contact target
			if c != (model.Contact{}) {
				contacts, err := SendFindContactMessage(&c, net.GetIdentity(), target, k)
				//should this check whether target is you?
				if err != nil {
					if contacts != nil {
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
	for i := range closestContacts {
		contactIn <- closestContacts[i]
		go run(contactIn, contactOut)
	}

	numWorkers := len(closestContacts)
	for numWorkers > 0 {
		receivedContact := <-contactOut
		// If closer than one of closestContacts, insert it instead and add it to contactIn
		// If not insert an empty contact to kill a worker
		foundCloser := false

		for i, contact := range closestContacts {
			if receivedContact.Less(&contact) { // Check if closer
				closestContacts[i] = receivedContact

				contactIn <- receivedContact // Queue up another contact for contacting
				foundCloser = true
				break
			}
		}

		if !foundCloser {
			contactIn <- model.Contact{} // Kill a worker if it couldn't find any closer contacts
			numWorkers--
		}
	}

	return closestContacts
}

func LookupData(net *model.KademliaNetwork, fileID *model.KademliaID) ([]byte, error) {

	//check if present locally
	data, found := net.GetData(fileID)
	if found {
		return data, nil
	}

	closestContacts := net.GetContacts(fileID, parallelism)
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
				data, contacts, err := SendFindDataMessage(&c, net.GetIdentity(), fileID, 3) // Best value for nbNeighbors?
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
	for _, c := range closestContacts {
		contactIn <- c
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

	return nil, errors.New("file could not be found ")
}

func StoreData(net *model.KademliaNetwork, data []byte) (fileID model.KademliaID, err error) {
	targetID := model.NewKademliaID(data)
	contacts := append(LookupContact(net, targetID), *net.GetIdentity())

	//fmt.Print("ID is '%s'", targetID.String())

	for _, contact := range contacts {
		err = SendStoreMessage(&contact, net.GetIdentity(), data)
		if err != nil {
			return fileID, err
		}
	}

	return *targetID, nil
}

func JoinNetwork(net *model.KademliaNetwork, IP string) error {
	target := model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: IP,
	}

	foundContacts, err := SendFindContactMessage(&target, net.GetIdentity(), net.GetIdentity().ID, k)
	if err != nil {
		return err
	}

	for _, contact := range foundContacts {
		net.RegisterContact(contact)
	}

	return nil
}
