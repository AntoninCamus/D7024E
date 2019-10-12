package model

import (
	"container/list"
)

// contactBucket definition
// contains a List
type contactBucket struct {
	list *list.List
}

// newBucket returns a new instance of a contactBucket
func newBucket() *contactBucket {
	b := &contactBucket{}
	b.list = list.New()
	return b
}

// addContact adds the Contact to the front of the contactBucket
// or moves it to the front of the contactBucket if it already existed
func (bucket *contactBucket) addContact(contact Contact) {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushFront(contact)
		}
	} else {
		bucket.list.MoveToFront(element)
	}
}

// getContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *contactBucket) getContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

// len return the maxSize of the contactBucket
func (bucket *contactBucket) len() int {
	return bucket.list.Len()
}
