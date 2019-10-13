package kademlia

import (
	"gotest.tools/assert"
	"os"
	"testing"

	"github.com/LHJ/D7024E/kademlia/model"
	"os/exec"
)

func Test_userAPI(t *testing.T) {
	tk := model.NewKademliaNetwork(model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "127.0.0.1",
	})
	sigChan := make(chan os.Signal, 1)

	StartGrpcServer(tk)
	StartRestServer(tk, sigChan)

	// Store
	cmd := exec.Command("../client.py", "store", "--file=TEST", "localhost:8080")
	stdOut, _ := cmd.CombinedOutput()

	assert.Equal(t, string(stdOut)[11:23], "Status:  200")

	// Find
	id := model.NewKademliaID([]byte("TEST"))
	cmd = exec.Command("../client.py", "find", "--id="+id.String(), "localhost:8080")
	stdOut, _ = cmd.CombinedOutput()

	assert.Equal(t, string(stdOut)[11:23], "Status:  200")

	id = model.NewKademliaID([]byte("MISSING FILE"))
	cmd = exec.Command("../client.py", "find", "--id="+id.String(), "localhost:8080")
	stdOut, _ = cmd.CombinedOutput()

	assert.Equal(t, string(stdOut)[11:23], "Status:  404")

	// Exit
	cmd = exec.Command("../client.py", "exit", "localhost:8080")
	stdOut, _ = cmd.CombinedOutput()

	assert.Equal(t, string(stdOut)[11:23], "Status:  200")
}
