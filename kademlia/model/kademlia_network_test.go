package model

import (
	"testing"

	"gotest.tools/assert"
)

// LOCAL (THREAD SAFE, BASIC) FUNCTIONS :
func TestKademlia_GetIdentity(t *testing.T) {
	me := newContact(NewRandomKademliaID(), "127.0.0.1")

	kad := NewKademliaNetwork(me)
	assert.Equal(t, kad.GetIdentity().ID, me.ID)
	assert.Equal(t, kad.GetIdentity().Address, me.Address)
}

func TestKademlia_SaveAndGetData(t *testing.T) {
	me := newContact(NewRandomKademliaID(), "127.0.0.1")
	kad := NewKademliaNetwork(me)

	fileID := NewRandomKademliaID()
	content := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit.")
	err := kad.SaveData(fileID, content)
	assert.NilError(t, err)

	contentFound, found := kad.GetData(fileID)
	assert.Assert(t, found)
	assert.DeepEqual(t, contentFound, content)

	contentFound, found = kad.GetData(NewRandomKademliaID())
	assert.Assert(t, !found)
	assert.Equal(t, len(contentFound), 0)
}

func TestKademlia_RegisterAndGetContact(t *testing.T) {
	me := newContact(NewRandomKademliaID(), "127.0.0.1")
	kad := NewKademliaNetwork(me)

	emptyContacts := kad.GetContacts(me.ID, 10)
	assert.Equal(t, len(emptyContacts), 0)

	c1 := newContact(NewRandomKademliaID(), "127.0.1.1")
	c2 := newContact(NewRandomKademliaID(), "127.0.2.1")
	kad.RegisterContact(&c1)
	kad.RegisterContact(&c2)
	contacts := kad.GetContacts(me.ID, 10)
	assert.Equal(t, len(contacts), 2)

	kad.RegisterContact(&c2)
	contacts = kad.GetContacts(me.ID, 10)
	assert.Equal(t, len(contacts), 2)

	contains := func(arr []Contact, elem Contact) bool {
		for _, a := range arr {
			if a.ID == elem.ID && a.Address == elem.Address {
				return true
			}
		}
		return false
	}
	assert.Assert(t, contains(contacts, c1))
	assert.Assert(t, contains(contacts, c2))
}
