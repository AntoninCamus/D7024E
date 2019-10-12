package kademlia

import (
	"gotest.tools/assert"
	"testing"
)

func TestHardware_GetContactFromHW_ShouldNotChange(t *testing.T) {
	c := GetContactFromHW()
	c2 := GetContactFromHW()
	assert.Equal(t, c.ID.String(), c2.ID.String())
	assert.Equal(t, c.Address, c2.Address)
}
