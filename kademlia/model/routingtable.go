package model

import (
	"fmt"
)

//BucketSize is the size of one bucket
const BucketSize = 3

// routingTable definition
// keeps a refrence contact of me and an array of buckets
type routingTable struct {
	me      Contact
	buckets [IDLength * 8]*contactBucket
}

// newRoutingTable returns a new instance of a RoutingTable
func newRoutingTable(me Contact) *routingTable {
	table := &routingTable{}
	for i := 0; i < IDLength*8; i++ {
		table.buckets[i] = newBucket()
	}
	table.me = me
	return table
}

// addContact add a new contact to the correct Bucket
func (routingTable *routingTable) addContact(contact Contact) {
	bucketIndex := routingTable.getBucketIndex(contact.ID)
	b := routingTable.buckets[bucketIndex]
	b.addContact(contact)
}

// containContact return true if a bucket contain this contact already
func (routingTable *routingTable) containContact(id KademliaID) bool {
	idx := routingTable.getBucketIndex(&id)
	for e := routingTable.buckets[idx].list.Front(); e != nil; e = e.Next() {
		s, ok := e.Value.(Contact)
		if ok && *s.ID == id {
			return true
		}
	}
	return false
}

// findClosestContacts finds the count closest Contacts to the target in the RoutingTable
func (routingTable *routingTable) findClosestContacts(target *KademliaID, count int) []Contact {
	var candidates contactCandidates
	bucketIndex := routingTable.getBucketIndex(target)
	b := routingTable.buckets[bucketIndex]

	candidates.append(b.getContactAndCalcDistance(target))

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < IDLength*8) && candidates.Len() < count; i++ {
		if bucketIndex-i >= 0 {
			b = routingTable.buckets[bucketIndex-i]
			candidates.append(b.getContactAndCalcDistance(target))
		}
		if bucketIndex+i < IDLength*8 {
			b = routingTable.buckets[bucketIndex+i]
			candidates.append(b.getContactAndCalcDistance(target))
		}
	}

	candidates.sort()

	if count > candidates.Len() {
		count = candidates.Len()
	}

	return candidates.getContacts(count)
}

// getBucketIndex get the correct Bucket index for the KademliaID
func (routingTable *routingTable) getBucketIndex(id *KademliaID) int {
	distance := id.calcDistance(routingTable.me.ID)
	for i := 0; i < IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}

	return IDLength*8 - 1
}

// getMe return current identity
func (routingTable *routingTable) getMe() Contact {
	return routingTable.me
}

func (routingTable *routingTable) String() string {
	var skip = false
	var ret = "["
	for _, val := range routingTable.buckets {
		if val.len() > 0 {
			ret += fmt.Sprintf("%s,", val.string())
			skip = false
		} else {
			if !skip {
				ret += "..."
				skip = true
			}
		}
	}
	ret += "]"
	return ret
}
