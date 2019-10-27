package kademlia

import (
	"testing"

	"github.com/LHJ/D7024E/kademlia/model"
	"gotest.tools/assert"
)

func Test_Ping(t *testing.T) {
	tk := model.NewKademliaNetwork(model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "127.0.0.1",
	})
	s := StartGrpcServer(tk)

	me := model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "127.0.0.1",
	}

	// Send a normal ping that should work
	p_normal := sendPingMessage(&me, false)
	assert.Assert(t, p_normal)
	s.GracefulStop()

	// Send a ping on a offline node
	p_offline := sendPingMessage(&me, false)
	assert.Assert(t, !p_offline)
}

func Test_FindContact(t *testing.T) {
	tk := model.NewKademliaNetwork(model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "127.0.0.1",
	})
	s := StartGrpcServer(tk)

	me := model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "127.0.0.1",
	}

	_, err := sendFindContactMessage(&me, &me, model.NewRandomKademliaID(), 0)
	assert.NilError(t, err)

	s.GracefulStop()

	_, err = sendFindContactMessage(&me, &me, model.NewRandomKademliaID(), 0)
	assert.Error(t, err, "context deadline exceeded")
}

func Test_StoreDataCall(t *testing.T) {
	tk := model.NewKademliaNetwork(model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "127.0.0.1",
	})
	s := StartGrpcServer(tk)

	me := model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "127.0.0.1",
	}

	err := sendStoreMessage(&me, &me, []byte("TEST1"))
	assert.NilError(t, err)

	s.GracefulStop()

	err = sendStoreMessage(&me, &me, []byte("TEST2"))
	assert.Error(t, err, "context deadline exceeded")
}

func Test_FindDataCall(t *testing.T) {
	tk := model.NewKademliaNetwork(model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "127.0.0.1",
	})
	s := StartGrpcServer(tk)

	me := tk.GetIdentity()
	other := model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "127.0.0.1",
	}

	data := []byte("TEST")
	id := model.NewKademliaID(data)

	err := sendStoreMessage(&me, &other, data)
	assert.NilError(t, err)

	dataReceived, _, err := sendFindDataMessage(&me, &other, id, 1)
	assert.NilError(t, err)
	assert.Equal(t, string(dataReceived), string(data))

	dataReceived, contacts, err := sendFindDataMessage(&me, &other, model.NewRandomKademliaID(), 1)
	assert.NilError(t, err)
	assert.Equal(t, len(dataReceived), 0)
	assert.Equal(t, len(contacts), 1)
	assert.Equal(t, contacts[0].ID.String(), other.ID.String())
	s.GracefulStop()
}
