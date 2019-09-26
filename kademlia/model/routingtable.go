package model

const bucketSize = 20

// RoutingTable definition
// keeps a refrence contact of Me and an array of buckets
type RoutingTable struct {
	Me      Contact
	buckets [IDLength * 8]*bucket
}

// NewRoutingTable returns a new instance of a RoutingTable
func NewRoutingTable(me Contact) *RoutingTable {
	routingTable := &RoutingTable{}
	for i := 0; i < IDLength*8; i++ {
		routingTable.buckets[i] = newBucket()
	}
	routingTable.Me = me
	return routingTable
}

// AddContact add a new contact to the correct Bucket
func (routingTable *RoutingTable) AddContact(contact Contact) {
	bucketIndex := routingTable.getBucketIndex(contact.ID)
	b := routingTable.buckets[bucketIndex]
	b.addContact(contact)
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
	distance := id.calcDistance(routingTable.Me.ID)
	for i := 0; i < IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}

	return IDLength*8 - 1
}
