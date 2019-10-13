package model

import (
	"fmt"
)

const bucketSize = 3

// RoutingTable definition
// keeps a refrence contact of me and an array of buckets
type RoutingTable struct {
	me      Contact
	buckets [IDLength * 8]*contactBucket
}

// NewRoutingTable returns a new instance of a RoutingTable
func NewRoutingTable(me Contact) *RoutingTable {
	routingTable := &RoutingTable{}
	for i := 0; i < IDLength*8; i++ {
		routingTable.buckets[i] = newBucket()
	}
	routingTable.me = me
	return routingTable
}

// AddContact add a new contact to the correct Bucket
func (routingTable *RoutingTable) AddContact(contact Contact) {
	bucketIndex := routingTable.getBucketIndex(contact.ID)
	b := routingTable.buckets[bucketIndex]
	b.addContact(contact)
}

func (routingTable *RoutingTable) ContainContact(id KademliaID) bool {
	idx := routingTable.getBucketIndex(&id)
	for e := routingTable.buckets[idx].list.Front(); e != nil; e = e.Next() {
		s, ok := e.Value.(Contact)
		if ok && *s.ID == id {
			return true
		}
	}
	return false
}

// FindClosestContacts finds the count closest Contacts to the target in the RoutingTable
func (routingTable *RoutingTable) FindClosestContacts(target *KademliaID, count int) []Contact {
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
func (routingTable *RoutingTable) getBucketIndex(id *KademliaID) int {
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

func (routingTable *RoutingTable) GetMe() Contact {
	return routingTable.me
}

func (routingTable *RoutingTable) String() string {
	var skip = false
	var ret = "["
	for _, val := range routingTable.buckets {
		if val.len() > 0 {
			ret += fmt.Sprintf("%s,", val.String())
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
