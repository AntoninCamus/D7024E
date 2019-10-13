package model

import (
	"fmt"
	"gotest.tools/assert"
	"sort"
	"testing"
)

func TestRoutingTable_ShouldReturnOrdered(t *testing.T) {

	var searchedID = "2111111400000000000000000000000000000000"

	var ids = []string{
		"ffffffff00000000000000000000000000000000",
		"1111111400000000000000000000000000000000",
		"1111111100000000000000000000000000000000",
		"1111111200000000000000000000000000000000",
		"1111111300000000000000000000000000000000",
		searchedID,
	}

	rt := newRoutingTable(newContact(KademliaIDFromString(ids[0]), "localhost:8000"))

	for i, c := range ids {
		rt.addContact(newContact(KademliaIDFromString(c), fmt.Sprintf("localhost:800%d", i+1)))
	}

	contacts := rt.findClosestContacts(KademliaIDFromString(searchedID), 20)
	sort.Slice(ids, func(i, j int) bool {
		ci := newContact(KademliaIDFromString(ids[i]), "localhost:8000")
		cj := newContact(KademliaIDFromString(ids[j]), "localhost:8000")
		ci.CalcDistance(KademliaIDFromString(searchedID))
		cj.CalcDistance(KademliaIDFromString(searchedID))
		return ci.less(&cj)
	})

	for i := range contacts {
		t.Log(contacts[i].String())
		assert.Equal(t, contacts[i].ID.String(), ids[i])
	}
}
