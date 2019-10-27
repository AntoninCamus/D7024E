package model

import (
	"testing"
)

func TestNewBucket(t *testing.T) {
	bucket1 := newBucket()
	bucket2 := newBucket()
	if bucket1 == bucket2 {
		t.Errorf("The two new buckets are the same object")
	}
}
func TestAddContact(t *testing.T) {
	testBucket := newBucket()

	testContact1 := newContact(NewRandomKademliaID(), "someaddress")
	testBucket.addContact(testContact1)

	testContact2 := newContact(NewRandomKademliaID(), "someotheraddress")
	testBucket.addContact(testContact2)

	if testBucket.list.Front().Value.(Contact).ID != testContact1.ID && testBucket.list.Back().Value.(Contact).ID != testContact1.ID {
		//log.Print(testBucket.list.Front().Value.(Contact).ID) first not found anywhere
		//log.Print(testBucket.list.Back().Value.(Contact).ID)
		t.Errorf("addContact didn't add the contacts with id %s to the list, or it has some other values", testContact1.ID.String())
	} else if testBucket.list.Front().Value.(Contact).ID != testContact2.ID && testBucket.list.Back().Value.(Contact) != testContact2 {
		t.Errorf("addContact didn't add the contact with id %s to the list, or it has some other values too", testContact2.ID.String())
	}

	testBucket.addContact(testContact1)
	if testBucket.list.Front().Value != testContact1 {
		t.Errorf("testContact1 was not moved to the front")
	}
	//for i := 0; testBucket.len() < BucketSize; i++{
	//	testBucket.addContact(newContact(NewRandomKademliaID(), "address " + string(i)))
	//}

}

func TestRemoveContact(t *testing.T) {
	testBucket := newBucket()
	testContact := newContact(NewRandomKademliaID(), "someaddress")
	initialLen := testBucket.len()
	testBucket.addContact(testContact)
	addedLen := testBucket.len()
	if addedLen <= initialLen {
		t.Errorf("Bucket did not increase in size")
	}
	testBucket.removeContact(testContact)
	if addedLen < testBucket.len() {
		t.Errorf("Bucket did not decrease in size")
	}
}
func TestGetContactAndCalcDistance(t *testing.T) {
	testBucket := newBucket()
	testContact := newContact(NewRandomKademliaID(), "someaddress")
	testBucket.addContact(testContact)

	targetContact := newContact(NewRandomKademliaID(), "targetAddress")

	rangedContacts := testBucket.getContactAndCalcDistance(targetContact.ID)

	//really, is it the job of this test to check that the CalcDistance does it's job?

	if len(rangedContacts) != testBucket.len() {
		t.Errorf("getContactAndCalcDistance returned an array of contacts of different size than the buckets")
	}

}

func TestLen(t *testing.T) {
	testBucket := newBucket()
	testContact := newContact(NewRandomKademliaID(), "someaddress")
	initialLen := testBucket.len()
	testBucket.addContact(testContact)
	addedLen := testBucket.len()
	if addedLen <= initialLen {
		t.Errorf("Bucket did not increase in size")
	}
}
