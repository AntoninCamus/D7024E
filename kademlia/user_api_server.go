package kademlia

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/LHJ/D7024E/kademlia/model"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"syscall"
)

type messageAnswer struct {
	Message string `json:"message"`
}

type restService struct {
	sigChannel      chan os.Signal
	kademliaNetwork *model.KademliaNetwork
}

// StartRestServer start the REST User API
func StartRestServer(k *model.KademliaNetwork, s chan os.Signal) *http.Server {
	fmt.Println("Starting server...")

	service := restService{
		sigChannel:      s,
		kademliaNetwork: k,
	}

	srv := http.Server{Addr: ":8080", Handler: nil}

	http.HandleFunc("/kademlia/file", service.findstore)
	http.HandleFunc("/node/exit", service.exitServer)

	serving := func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
	go serving()

	return &srv
}

func (s *restService) findstore(w http.ResponseWriter, r *http.Request) {
	var answer interface{} // Generic

	if r.Method == "POST" { // Store
		fileID, err := store(w, r, s.kademliaNetwork)
		if err != nil {
			log.Printf("API ERROR : %s", err.Error())
			return
		}
		answer = struct{ fileID string }{fileID: fileID}
	} else if r.Method == "GET" { // Find
		data, err := find(w, r, s.kademliaNetwork)
		if err != nil {
			log.Printf("API ERROR : %s", err.Error())
			return
		}
		answer = struct{ data string }{data: data}
	}

	// Write response
	js, err := json.Marshal(answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func find(w http.ResponseWriter, r *http.Request, network *model.KademliaNetwork) (string, error) {
	// Read input
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return "", err
	}
	id := r.Form.Get("id")

	// Retrieve file
	kademliaID := model.KademliaIDFromString(id)
	file, err := LookupData(network, kademliaID)
	if err != nil {
		http.Error(w, "Error while retrieving file", http.StatusNotFound)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(file), nil
}

func store(w http.ResponseWriter, r *http.Request, network *model.KademliaNetwork) (string, error) {
	// Read input
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return "", nil
	}

	// Store file
	file := []byte(string(body))
	id, err := StoreData(network, file)
	if err != nil {
		http.Error(w, "Error while storing file", http.StatusInternalServerError)
		return "", err
	}
	log.Print("Store successful, new state is :", network.PrintFileState())
	return id.String(), nil
}

func (s *restService) exitServer(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(
		messageAnswer{Message: "Server shutting down..."},
	)

	defer func() {
		s.sigChannel <- syscall.SIGTERM
	}()
}
