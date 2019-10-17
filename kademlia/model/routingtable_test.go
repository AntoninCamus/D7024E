package model

import (
	"fmt"
	"gotest.tools/assert"
	"testing"
)

func TestRoutingTable_ShouldReturnOrdered(t *testing.T) {

	var nodeID = KademliaIDFromString("ffffffff00000000000000000000000000000000")
	var searchedID = KademliaIDFromString("2111111400000000000000000000000000000000")

	var ids = []*KademliaID{
		KademliaIDFromString("1111111100000000000000000000000000000000"),
		searchedID,
		KademliaIDFromString("1111111200000000000000000000000000000000"),
		KademliaIDFromString("1111111300000000000000000000000000000000"),
		nodeID,
		KademliaIDFromString("1111111400000000000000000000000000000000"),
	}

	var expectedtIds = []*KademliaID{
		searchedID,
		KademliaIDFromString("1111111400000000000000000000000000000000"),
		KademliaIDFromString("1111111100000000000000000000000000000000"),
		KademliaIDFromString("1111111200000000000000000000000000000000"),
		KademliaIDFromString("1111111300000000000000000000000000000000"),
		nodeID,
	}

	for i := range expectedtIds[:len(expectedtIds)-1] {
		assert.Assert(t, expectedtIds[i].calcDistance(searchedID).less(expectedtIds[i+1].calcDistance(searchedID)))
	}

	rt := newRoutingTable(newContact(KademliaIDFromString("ffffffff00000000000000000000000000000000"), "localhost:8000"))

	for i, c := range ids {
		c := newContact(c, fmt.Sprintf("localhost:800%d", i+1))
		rt.addContact(c)
	}

	contacts := rt.findClosestContacts(searchedID, 20)

	for i := range contacts {
		assert.Assert(t, contacts[i].ID.equals(expectedtIds[i]))
	}
}
