package kademlia

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/LHJ/D7024E/kademlia/model"
	pb "github.com/LHJ/D7024E/protogen"
	"google.golang.org/grpc"
)

func connect(address string) (pb.InternalApiServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", address, grpcPort),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(3*time.Second),
	)
	if err != nil {
		return nil, nil, err
	}

	client := pb.NewInternalApiServiceClient(conn)
	return client, conn, nil
}

// sendPingMessage ping the provided contact and return if it is present or not.
func sendPingMessage(target *model.Contact) bool {
	// Open gRPC connection
	client, conn, err := connect(target.Address)
	if err != nil {
		log.Print(err)
		return false
	}
	defer func() {
		if err = conn.Close(); err != nil {
			log.Print(err)
		}
	}()

	ans, err := client.PingCall(
		context.Background(),
		&pb.PingRequest{SenderKademliaId: model.NewRandomKademliaID()[:]},
	)
	if err != nil {
		log.Print(err)
		return false
	}

	return len(ans.GetReceiverKademliaId()) == model.IDLength
}

// sendFindContactMessage ask to the provided node for the nbNeighbors closest neighbors of the searchedContactID provided, and returns them.
func sendFindContactMessage(target *model.Contact, me *model.Contact, searchedContactID *model.KademliaID, nbNeighbors int) (contacts []*model.Contact, err error) {
	// Open gRPC connection
	client, conn, err := connect(target.Address)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = conn.Close(); err != nil {
			log.Print(err)
		}
	}()

	ans, err := client.FindContactCall(
		context.Background(),
		&pb.FindContactRequest{
			Src: &pb.Contact{
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

// sendFindDataMessage ask to the provided node for the file identified by the provided fileID, and returns it.
// If data was not found it act as SendFindContactMessage.
func sendFindDataMessage(target *model.Contact, me *model.Contact, searchedFileID *model.KademliaID, nbNeighbors int) ([]byte, []*model.Contact, error) {
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
		&pb.FindDataRequest{
			Src: &pb.Contact{
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
	case *pb.FindDataAnswer_DataFound:
		return ans.GetDataFound(), nil, nil

	case *pb.FindDataAnswer_DataNotFound:
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

// sendStoreMessage ask to the provided node to store the file, and returns the corresponding ID.
func sendStoreMessage(target *model.Contact, me *model.Contact, data []byte) error {
	if target.Address == "" {
		return fmt.Errorf("target is invalid %s", target.String())
	}

	client, conn, err := connect(target.Address)
	if err != nil {
		log.Print("Unable to connect to", target.Address)
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Print(err)
		}
	}()

	done, err := client.StoreDataCall(
		context.Background(),
		&pb.StoreDataRequest{
			Src: &pb.Contact{
				ID:      me.ID[:],
				Address: me.Address,
			},
			Data: data,
		},
	)
	if err != nil {
		return fmt.Errorf("unable store onto %s, got error '%s'", target.Address, err.Error())
	}
	if !done.Ok {
		return errors.New("distant node was unable to store data")
	}
	return nil
}
