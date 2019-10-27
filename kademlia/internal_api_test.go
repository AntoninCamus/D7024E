package kademlia

import (
	"fmt"
	"testing"

	"github.com/LHJ/D7024E/kademlia/model"
	"gotest.tools/assert"
)

func Test_Ping(t *testing.T) {
	network := model.NewKademliaNetwork(model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "127.0.0.1",
	})

	s := StartGrpcServer(network)
	source := model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "127.0.0.1",
	}

	// Should work
	answerValid := sendPingMessage(&source, false)
	assert.Assert(t, answerValid)
	s.GracefulStop()

	// Should timeout
	s.GracefulStop()
	answerTimout := sendPingMessage(&source, false)
	assert.Assert(t, !answerTimout)
}

func Test_FindContact(t *testing.T) {
	net := model.NewKademliaNetwork(model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "127.0.0.1",
	})
	s := StartGrpcServer(net)

	source := model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "127.0.0.1",
	}
	target := net.GetIdentity()

	// Should work
	contacts, err := sendFindContactMessage(
		&target,
		&source,
		model.NewRandomKademliaID(),
		5,
	)
	assert.NilError(t, err)
	assert.Equal(t, contacts[0].ID.String(), source.ID.String())

	// Should fail if no address
	invalidTarget := model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "",
	}
	_, err = sendFindContactMessage(&invalidTarget, &source, model.NewRandomKademliaID(), 5)
	assert.Error(t, err, fmt.Sprintf("target is invalid %s", invalidTarget.String()))

	// Should timeout
	s.GracefulStop()
	_, err = sendFindContactMessage(&target, &target, model.NewRandomKademliaID(), 5)
	assert.Error(t, err, "context deadline exceeded")
}

func Test_StoreDataCall(t *testing.T) {
	net := model.NewKademliaNetwork(model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "127.0.0.1",
	})
	s := StartGrpcServer(net)

	source := model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "127.0.0.1",
	}
	target := net.GetIdentity()

	// Should work
	err := sendStoreMessage(&target, &source, []byte("TEST1"))
	assert.NilError(t, err)

	// Should fail if no address
	invalidTarget := model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "",
	}
	err = sendStoreMessage(&invalidTarget, &source, []byte("TEST1"))
	assert.Error(t, err, fmt.Sprintf("target is invalid %s", invalidTarget.String()))

	// Should timeout
	s.GracefulStop()
	err = sendStoreMessage(&target, &source, []byte("TEST2"))
	assert.Error(t, err, "context deadline exceeded")
}

func Test_FindDataCall(t *testing.T) {
	net := model.NewKademliaNetwork(model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "127.0.0.1",
	})
	s := StartGrpcServer(net)

	source := model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "127.0.0.1",
	}
	target := net.GetIdentity()

	data := []byte("TEST")
	id := model.NewKademliaID(data)

	// Should work
	err := sendStoreMessage(&target, &source, data)
	assert.NilError(t, err)

	dataReceived, contacts, err := sendFindDataMessage(&target, &source, id, 1)
	assert.NilError(t, err)
	assert.Equal(t, string(dataReceived), string(data))
	assert.Equal(t, len(contacts), 0)

	dataReceived, contacts, err = sendFindDataMessage(&target, &source, model.NewRandomKademliaID(), 1)
	assert.NilError(t, err)
	assert.Equal(t, len(dataReceived), 0)
	assert.Equal(t, contacts[0].ID.String(), source.ID.String())

	// Should fail if no address
	invalidTarget := model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "",
	}
	_, _, err = sendFindDataMessage(&invalidTarget, &source, id, 1)
	assert.Error(t, err, fmt.Sprintf("target is invalid %s", invalidTarget.String()))

	// Should timeout
	s.GracefulStop()
	_, _, err = sendFindDataMessage(&target, &source, id, 1)
	assert.Error(t, err, "context deadline exceeded")
}
