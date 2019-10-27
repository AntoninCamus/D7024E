package kademlia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/LHJ/D7024E/kademlia/model"
	"gotest.tools/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestUserApi(t *testing.T) {
	testTable := []func(t *testing.T, net *model.KademliaNetwork, sigChan chan os.Signal){findFile, postFile, exitNode}

	net := model.NewKademliaNetwork(model.Contact{
		ID:      model.NewRandomKademliaID(),
		Address: "127.0.0.1",
	})
	sigChan := make(chan os.Signal, 1)

	g := StartGrpcServer(net)
	h := StartRestServer(net, sigChan)

	for _, test := range testTable {
		test(t, net, sigChan)
	}

	g.GracefulStop()
	h.Close()
}

func postFile(t *testing.T, net *model.KademliaNetwork, sigChan chan os.Signal) {

	fileToStore := "123"

	response, err := http.Post(
		"http://localhost:8080/kademlia/file",
		"application/raw",
		bytes.NewBufferString(fileToStore),
	)
	assert.NilError(t, err)
	assert.Equal(t, response.StatusCode, 200)

	data, err := ioutil.ReadAll(response.Body)
	assert.NilError(t, err)

	var answer StoreAnswer
	err = json.Unmarshal(data, &answer)
	assert.NilError(t, err)

	//answer.FileID should be a valid KademliaID
	assert.Equal(t, len(answer.FileID), len(model.NewRandomKademliaID().String()))
	id := model.KademliaIDFromString(answer.FileID)

	//file should be in the storage
	data, found := net.GetData(id)
	assert.Assert(t, found)
	assert.Equal(t, string(data), fileToStore)
}

func findFile(t *testing.T, net *model.KademliaNetwork, sigChan chan os.Signal) {
	fileToFind := "123"
	idToFind := model.NewRandomKademliaID()
	err := net.SaveData(idToFind, []byte(fileToFind))
	assert.NilError(t, err)

	// Find an existing file should return 200 and no error
	response, err := http.Get(fmt.Sprintf(
		"http://localhost:8080/kademlia/file?id=%s",
		idToFind.String(),
	))
	assert.NilError(t, err)
	assert.Equal(t, response.StatusCode, 200)

	// Then the answer should be equal to the ID inserted
	data, err := ioutil.ReadAll(response.Body)
	assert.NilError(t, err)
	var answer FindAnswer
	err = json.Unmarshal(data, &answer)
	assert.NilError(t, err)
	assert.Equal(t, answer.Data, fileToFind)

	// Finally seaching for a non existing one should return a 404
	response, err = http.Get(fmt.Sprintf(
		"http://localhost:8080/kademlia/file?id=%s",
		model.NewRandomKademliaID().String(),
	))
	assert.NilError(t, err)
	assert.Equal(t, response.StatusCode, 404)

}

func exitNode(t *testing.T, net *model.KademliaNetwork, sigChan chan os.Signal) {
	// Channel should be empty before calling exit
	assert.Equal(t, len(sigChan), 0)

	response, err := http.Get("http://localhost:8080/node/exit")
	assert.NilError(t, err)
	assert.Equal(t, response.StatusCode, 200)

	// Channel should contain one value after calling exit
	assert.Equal(t, len(sigChan), 1)
}
