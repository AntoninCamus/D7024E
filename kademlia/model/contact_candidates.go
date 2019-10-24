package model

import "sort"

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
		return candidates.contacts[i].Less(&candidates.contacts[j])
	})
}

// Len returns the length of the ContactCandidates
func (candidates *contactCandidates) Len() int {
	return len(candidates.contacts)
}
