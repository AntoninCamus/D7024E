package networking

import (
	"testing"

	"github.com/LHJ/D7024E/kademlia/model"
	"gotest.tools/assert"
)

func TestPing(t *testing.T) {
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
