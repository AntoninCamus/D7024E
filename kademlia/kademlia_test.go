package kademlia

import (
	"github.com/LHJ/D7024E/kademlia/model"
	"gotest.tools/assert"
	"testing"
)

// LOCAL (THREAD SAFE, BASIC) FUNCTIONS :
func TestKademlia_GetIdentity(t *testing.T) {
	me := model.NewContact(model.NewRandomKademliaID(), "127.0.0.1")

	kad := Init(me)
	assert.Equal(t, kad.GetIdentity(), me)
}

func TestKademlia_SaveAndGetData(t *testing.T) {
	me := model.NewContact(model.NewRandomKademliaID(), "127.0.0.1")
	kad := Init(me)

	fileID := model.NewRandomKademliaID()
	content := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit.")
	err := kad.SaveData(fileID, content)
	assert.NilError(t, err)

	contentFound, found := kad.GetData(fileID)
	assert.Assert(t, found)
	assert.DeepEqual(t, contentFound, content)

	contentFound, found = kad.GetData(model.NewRandomKademliaID())
	assert.Assert(t, !found)
	assert.Equal(t, len(contentFound), 0)
}

func TestKademlia_RegisterAndGetContact(t *testing.T) {
	me := model.NewContact(model.NewRandomKademliaID(), "127.0.0.1")
	kad := Init(me)

	emptyContacts := kad.GetContacts(me.ID, 10)
	assert.Equal(t, len(emptyContacts), 0)

	c1 := model.NewContact(model.NewRandomKademliaID(), "127.0.1.1")
	c2 := model.NewContact(model.NewRandomKademliaID(), "127.0.2.1")
	kad.RegisterContact(&c1)
	kad.RegisterContact(&c2)
	contacts := kad.GetContacts(me.ID, 10)
	assert.Equal(t, len(contacts), 2)

	contains := func(arr []model.Contact, elem model.Contact) bool {
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
