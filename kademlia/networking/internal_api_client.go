package networking

import (
	"context"
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
		grpc.WithTimeout(1*time.Second),
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
	defer conn.Close()

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

// SendFindContactMessage ask to the provided node for the nbNeighbors closest neighbors of the nodeID provided, and returns them.
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
