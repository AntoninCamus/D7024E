package model

import (
	"fmt"
	"sort"
)

// Contact definition
// stores the KademliaID, the ip address and the distance
type Contact struct {
	ID       *KademliaID
	Address  string
	distance *KademliaID
}

// newContact returns a new instance of a Contact
func newContact(id *KademliaID, address string) Contact {
	return Contact{id, address, nil}
}

// CalcDistance calculates the distance to the target and
// fills the contacts distance field
func (contact *Contact) CalcDistance(target *KademliaID) {
	contact.distance = contact.ID.calcDistance(target)
}

// less returns true if contact.distance < otherContact.distance
func (contact *Contact) less(otherContact *Contact) bool {
	return contact.distance.less(otherContact.distance)
}

// String returns a simple string representation of a Contact
func (contact *Contact) String() string {
	return fmt.Sprintf(`contact("%s", "%s")`, contact.ID, contact.Address)
}

// ContactCandidates definition
// stores an array of Contacts
type contactCandidates struct {
	contacts []Contact
}

// Append an array of Contacts to the ContactCandidates
func (candidates *contactCandidates) append(contacts []Contact) {
	candidates.contacts = append(candidates.contacts, contacts...)
}

// GetContacts returns the first count number of Contacts
func (candidates *contactCandidates) getContacts(count int) []Contact {
	return candidates.contacts[:count]
}

// Sort the Contacts in ContactCandidates
func (candidates *contactCandidates) sort() {
	sort.Slice(candidates.contacts, func(i, j int) bool {
		return candidates.contacts[i].less(&candidates.contacts[j])
	})
}

// Len returns the length of the ContactCandidates
func (candidates *contactCandidates) Len() int {
	return len(candidates.contacts)
}
