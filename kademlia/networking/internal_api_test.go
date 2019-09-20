package networking

import (
	"testing"

	"github.com/LHJ/D7024E/kademlia/model"
	"gotest.tools/assert"
)

func Test_Ping(t *testing.T) {
	s := StartGrpcServer()

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
	s := StartGrpcServer()

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
