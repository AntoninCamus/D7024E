package model

// ContactSorter keep only a certain number of contacts, and only keep the closest to the target
type contactSorter struct {
	target   KademliaID
	contacts []Contact
}

// NewSorter returns a new instance of a ContactSorter, with a maximum of *maxSize*
func NewSorter(target KademliaID, size int) *contactSorter {
	return &contactSorter{
		target:   target,
		contacts: make([]Contact, size),
	}
}

// InsertContact insert the *contact* into the sorter.
// If the sorter is full, keep only the closer contacts.
// Returns if the sorter state was changed or not.
func (s *contactSorter) InsertContact(contact Contact) bool {
	// We start by calculating the distance of the target
	contact.CalcDistance(&s.target)

	// At start, we consider that maybe the contact is further than every of existing contacts
	further := -1

	// Then we iterate
	for i, c := range s.contacts {
		if c.ID == nil {
			// If the contact is empty, replace it and return there
			further = i
			break
		} else if c.ID == contact.ID {
			// If the contact is already present, interrupt
			further = -1
			break
		} else {
			// If the current contact is further than the current furthest replace it and set the position
			if further != -1 && c.less(&s.contacts[further]) {
				// If further != -1 we check the contacts list
				further = i
			} else if further == -1 && c.less(&contact) {
				// Else it means that the worse is the new one
				further = i
			}
		}
	}

	if further != -1 {
		// The provided contact should be inserted
		s.contacts[further] = contact
	}

	return further != -1
}

// GetContacts return the full internal state of the sorter.
func (s *contactSorter) GetContacts() []Contact {
	return s.contacts
}
