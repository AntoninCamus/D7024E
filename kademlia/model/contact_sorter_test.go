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
	targetID := NewRandomKademliaID()
	s := NewSorter(*targetID, 3)

	c := newContact(NewRandomKademliaID(), "")
	// Inserting two time the same should return true then false
	assert.Assert(t, s.InsertContact(c))
	assert.Assert(t, !s.InsertContact(c))

	//Another case :
	c1 := newContact(KademliaIDFromString("9c7dbe89c24da341bb751281926ddc11dbb656f1"), "10.0.0.201")
	c2 := newContact(KademliaIDFromString("9c7dbe89c24da341bb751281926ddc11dbb656f1"), "10.0.0.201")
	c1.CalcDistance(NewRandomKademliaID())
	c2.CalcDistance(NewRandomKademliaID())
	assert.Assert(t, s.InsertContact(c1))
	assert.Assert(t, !s.InsertContact(c2))
}

func TestContactSorter_InsertContact_InsertOnlyCloserValues(t *testing.T) {
	targetID := NewRandomKademliaID()
	s := NewSorter(*targetID, 5)

	for i := 0; i < 100; i++ {
		added := newContact(NewRandomKademliaID(), "")
		added.CalcDistance(targetID)

		// Let's compute the expected result :
		contactsBefore := s.GetContacts()
		expectedResult := false
		for _, c := range contactsBefore {
			if c.ID == added.ID {
				break
			} else if c.ID == nil || added.less(&c) {
				expectedResult = true
				break
			}
		}

		// Verify expected results
		assert.Equal(t, s.InsertContact(added), expectedResult)
		if expectedResult == true {
			// Expected replaced contact position
			expectedPosition := -1

			for j, c := range contactsBefore {
				t.Log(contactsBefore)
				t.Log(added)
				if c.ID == nil {
					t.Logf("NIL : %d", j)
					expectedPosition = j
					break
				} else {
					if expectedPosition == -1 && added.less(&c) {
						t.Logf("new worse : %d", j)
						expectedPosition = j
					} else if expectedPosition != -1 && contactsBefore[expectedPosition].less(&c) {
						t.Logf("new worse : %d", j)
						expectedPosition = j
					}
				}
			}
			assert.Assert(t, s.GetContacts()[expectedPosition].ID.equals(added.ID))
		}
	}
}
