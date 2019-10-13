package model

import (
	"gotest.tools/assert"
	"testing"
)

func TestContactSorter_InsertContact_ReturnsTrueWhenEmpty(t *testing.T) {
	id := NewRandomKademliaID()
	s := NewSorter(*id, 3)

	// InsertContact should return true the 3 first times
	for i := 0; i < 3; i++ {
		res := s.InsertContact(newContact(NewRandomKademliaID(), ""))
		assert.Assert(t, res)
	}
}

func TestContactSorter_InsertContact_ReturnsFalseWhenAlreadyPresent(t *testing.T) {
	id := NewRandomKademliaID()
	s := NewSorter(*id, 3)

	c := newContact(NewRandomKademliaID(), "")
	assert.Assert(t, s.InsertContact(c))
	// Inserting two time the same should return false
	assert.Assert(t, !s.InsertContact(c))
}

func TestContactSorter_InsertContact_InsertOnlyGreaterValues(t *testing.T) {
	SIZE := 5
	id := NewRandomKademliaID()
	s := NewSorter(*id, SIZE)

	// Fill the array
	for i := 0; i < SIZE; i++ {
		s.InsertContact(newContact(NewRandomKademliaID(), ""))
	}

	// Then add new ones
	for i := 0; i < SIZE; i++ {
		contact := newContact(NewRandomKademliaID(), "")

		// Compute expected result
		contact.CalcDistance(id)
		contacts := s.GetContacts()
		for _, c := range contacts {
			c.CalcDistance(id)
		}

		// Expected bool return
		expectedRes := false
		for _, c := range contacts {
			if c.less(&contact) {
				expectedRes = true
				break
			}
		}

		// Expected replaced contact position
		expectedReplacedContact := -1
		for j, c := range contacts {
			if expectedReplacedContact == -1 && c.less(&contact) {
				expectedReplacedContact = j
			} else if expectedReplacedContact != -1 && c.less(&contacts[expectedReplacedContact]) {
				expectedReplacedContact = j
			}
		}

		// Verify expected results
		assert.Equal(t, s.InsertContact(contact), expectedRes)
		if expectedReplacedContact != -1 {
			assert.Equal(t, s.GetContacts()[expectedReplacedContact].ID, contact.ID)
		}
	}
}
