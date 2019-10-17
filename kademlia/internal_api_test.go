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
	p_normal := sendPingMessage(&me)
	assert.Assert(t, p_normal)
	s.GracefulStop()

	// Send a ping on a offline node
	p_offline := sendPingMessage(&me)
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
	StartGrpcServer(tk)

	me := model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "127.0.0.1",
	}

	data := []byte("TEST")
	id := model.NewKademliaID(data)

	err := sendStoreMessage(&me, &me, data)
	assert.NilError(t, err)

	_, _, err = sendFindDataMessage(&me, &me, id, 1)

	assert.NilError(t, err)
}
