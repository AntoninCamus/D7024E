package model

import "testing"

func TestNewContact(t *testing.T)  {
	testContact1 := newContact(NewRandomKademliaID(), "someaddress")
	testContact2 := newContact(NewRandomKademliaID(), "someaddress")

	if testContact1 == testContact2 {
		t.Errorf("The two new contacts are the same object")
	}
}
func TestContact_CalcDistance(t *testing.T) {
	testContact1 := newContact(NewRandomKademliaID(), "someaddress")
	testContact2 := newContact(NewRandomKademliaID(), "someaddress")
	//testContact3 := newContact(NewRandomKademliaID(), "someaddress")

	distance := testContact1.ID.calcDistance(testContact2.ID)
	testContact1.CalcDistance(testContact2.ID)
	distance2 := testContact1.distance

	//distance3 := testContact1.ID.calcDistance(testContact3.ID)

	if *distance != *distance2 {
		t.Errorf("Calculated distance %s is not assigned to the contact but instead %s", distance, distance2)
	}
}

func TestContactCandidates_Len_Append(t *testing.T) {
	testContact1 := newContact(NewRandomKademliaID(), "someaddress")
	contactArray1 := []Contact{testContact1}

	testContact2 := newContact(NewRandomKademliaID(), "someaddress")
	contactArray2 := []Contact{testContact2}

	contacts := contactCandidates{contactArray1}

	shortDistance := contacts.Len()
	contacts.append(contactArray2)
	longDistance := contacts.Len()

	if shortDistance >= longDistance {
		t.Errorf("Append did not append the contact, or Len did not return the correct length")
	}
}

func TestGetContacts(t *testing.T){
	testContact1 := newContact(NewRandomKademliaID(), "someaddress")
	testContact2 := newContact(NewRandomKademliaID(), "someaddress")
	testContact3 := newContact(NewRandomKademliaID(), "someaddress")
	contactArray1 := []Contact{testContact1}
	contactArray2 := []Contact{testContact2, testContact3}
	verificationArray := []Contact{testContact1, testContact2}

	contactCanditateArray := contactCandidates{contactArray1}
	contactCanditateArray.append(contactArray2)

	returnContacts := contactCanditateArray.getContacts(2)
	for index, contact := range returnContacts{
		if contact != verificationArray[index]{
			t.Errorf("getContacts did not return the same array as input")
		}
	}
}

func TestContact_String(t *testing.T) {
	testContact := newContact(NewRandomKademliaID(), "someaddress")
	string := testContact.String()
	testString := "contact(\"" + testContact.ID.String() + "\", \"" + testContact.Address +"\")"
	if string != testString{
		t.Errorf("String() did not return correctly formatted single contact")//. %s, %s", string, testString)
	}
}

func TestSortLess(t *testing.T){
	ref := KademliaIDFromString("ffffffffffffffffffffffffffffffffffffffff")
	top := KademliaIDFromString("fffffffffffffffffffffffffffffffffffffffe")
	mid := KademliaIDFromString("7777777777777777777777777777777777777777")
	bot := KademliaIDFromString("0000000000000000000000000000000000000000")

	refContact := newContact(ref, "ref")
	topContact := newContact(top, "top")
	midContact := newContact(mid, "mid")
	botContact := newContact(bot, "bot")

	topContact.CalcDistance(refContact.ID)
	midContact.CalcDistance(refContact.ID)
	botContact.CalcDistance(refContact.ID)

	contactsArray := []Contact{botContact, topContact, midContact}

	contacts := contactCandidates{contactsArray}

	contacts.sort()

	sortedContacts := contacts.getContacts(3)

	if sortedContacts[0] != topContact || sortedContacts[1] != midContact || sortedContacts[2] != botContact {
		t.Errorf("sort() did not return a correctly sorted array")
	}
}