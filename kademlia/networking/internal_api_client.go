package networking

import (
	"github.com/LHJ/D7024E/kademlia/model"
)

// SendPingMessage ping the provided contact and return if it is present or not
func SendPingMessage(target *model.Contact) bool {
	// TODO
	return true
}

// SendFindContactMessage ask to the provided node for the nbNeighbors closest neighbors of the nodeID provided
func SendFindContactMessage(target *model.Contact, nodeID *model.KademliaID, nbNeighbors int) []*model.Contact {
	// TODO
	return make([]*model.Contact, 0)
}

// SendFindDataMessage ask to the provided node for the file identified by the provided fileID, and returns it.
func SendFindDataMessage(target *model.Contact, fileID *model.KademliaID) []byte {
	// TODO
	return make([]byte, 0)
}

// SendStoreMessage ask to the provided node to store the file, and returns the corresponding ID.
func SendStoreMessage(target *model.Contact, data []byte) *model.KademliaID {
	// TODO
	return model.NewRandomKademliaID()
}
