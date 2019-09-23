package networking

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/LHJ/D7024E/kademlia/model"
	"google.golang.org/grpc"
)

func connect(address string) (InternalApiServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", address, GrpcPort),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(3*time.Second),
	)
	if err != nil {
		return nil, nil, err
	}

	client := NewInternalApiServiceClient(conn)
	return client, conn, nil
}

// SendPingMessage ping the provided contact and return if it is present or not.
func SendPingMessage(target *model.Contact) bool {
	// Open gRPC connection
	client, conn, err := connect(target.Address)
	if err != nil {
		log.Print(err)
		return false
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Print(err)
		}
	}()

	ans, err := client.PingCall(
		context.Background(),
		&PingRequest{SenderKademliaId: model.NewRandomKademliaID()[:]},
	)
	if err != nil {
		log.Print(err)
		return false
	}

	return len(ans.GetReceiverKademliaId()) == model.IDLength
}

// SendFindContactMessage ask to the provided node for the nbNeighbors closest neighbors of the searchedContactID provided, and returns them.
func SendFindContactMessage(target *model.Contact, me *model.Contact, searchedContactID *model.KademliaID, nbNeighbors int) (contacts []*model.Contact, err error) {
	// Open gRPC connection
	client, conn, err := connect(target.Address)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Print(err)
		}
	}()

	ans, err := client.FindContactCall(
		context.Background(),
		&FindContactRequest{
			Me: &Contact{
				ID:      me.ID[:],
				Address: me.Address,
			},
			SearchedContactId: searchedContactID[:],
			NbNeighbors:       int32(nbNeighbors),
		},
	)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var modelContacts []*model.Contact
	for _, c := range ans.Contacts[:] {
		tmpID, err := model.KademliaIDFromBytes(c.ID)
		if err != nil {
			return nil, err
		}

		modelContacts = append(modelContacts, &model.Contact{
			ID:      tmpID,
			Address: c.Address,
		})
	}

	return modelContacts, nil
}

// SendFindDataMessage ask to the provided node for the file identified by the provided fileID, and returns it.
// If data was not found it act as SendFindContactMessage.
func SendFindDataMessage(target *model.Contact, me *model.Contact, searchedFileID *model.KademliaID, nbNeighbors int) ([]byte, []*model.Contact, error) {
	// Open gRPC connection
	client, conn, err := connect(target.Address)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Print(err)
		}
	}()

	ans, err := client.FindDataCall(
		context.Background(),
		&FindDataRequest{
			Me: &Contact{
				ID:      me.ID[:],
				Address: me.Address,
			},
			SearchedFileId: searchedFileID[:],
			NbNeighbors:    int32(nbNeighbors),
		},
	)
	if err != nil {
		log.Print(err)
		return nil, nil, err
	}

	switch ans.GetAnswer().(type) {
	case *FindDataAnswer_DataFound:
		return ans.GetDataFound(), nil, nil

	case *FindDataAnswer_DataNotFound:
		var modelContacts []*model.Contact
		for _, c := range ans.GetDataNotFound().Contacts[:] {
			tmpID, err := model.KademliaIDFromBytes(c.ID)
			if err != nil {
				return nil, nil, err
			}

			modelContacts = append(modelContacts, &model.Contact{
				ID:      tmpID,
				Address: c.Address,
			})
		}

		return nil, modelContacts, nil
	default:
		return nil, nil, errors.New("invalid answer content (neither found nor not found)")

	}
}

// SendStoreMessage ask to the provided node to store the file, and returns the corresponding ID.
func SendStoreMessage(target *model.Contact, me *model.Contact, data []byte) *model.KademliaID {
	// TODO
	return model.NewRandomKademliaID()
}
