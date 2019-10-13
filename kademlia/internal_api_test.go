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

	p1 := SendPingMessage(&me)
	assert.Equal(t, p1, true)
	s.GracefulStop()

	p2 := SendPingMessage(&me)
	assert.Equal(t, p2, false)
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

	_, err := SendFindContactMessage(&me, &me, model.NewRandomKademliaID(), 0)
	assert.NilError(t, err)

	s.GracefulStop()

	_, err = SendFindContactMessage(&me, &me, model.NewRandomKademliaID(), 0)
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

	err := SendStoreMessage(&me, &me,[]byte("TEST1"))
	assert.NilError(t, err)

	s.GracefulStop()

	err = SendStoreMessage(&me, &me,[]byte("TEST2"))
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

	SendStoreMessage(&me, &me,data)
	_, _, err := SendFindDataMessage(&me, &me, id,1)

	assert.NilError(t, err)
}