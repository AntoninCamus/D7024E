package kademlia

import (
	"encoding/json"
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
	log.Println("REST server is ready")

	return &srv
}

type StoreAnswer struct {
	FileID string
}

type FindAnswer struct {
	Data string
}

func (s *restService) findstore(w http.ResponseWriter, r *http.Request) {
	var jsonAnswer []byte // Generic

	if r.Method == "POST" { // Store
		fileID, err := store(w, r, s.kademliaNetwork)
		if err != nil {
			log.Printf("API error on POST during store : %s", err.Error())
			return
		}
		jsonAnswer, err = json.Marshal(StoreAnswer{FileID: fileID})
		if err != nil {
			log.Printf("API error on POST during marshal : %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "GET" { // Find
		data, err := find(w, r, s.kademliaNetwork)
		if err != nil {
			log.Printf("API error on GET during find : %s", err.Error())
			return
		}
		jsonAnswer, err = json.Marshal(FindAnswer{Data: data})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(jsonAnswer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func find(w http.ResponseWriter, r *http.Request, network *model.KademliaNetwork) (string, error) {
	// Read input
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "error reading request body", http.StatusBadRequest)
		return "", err
	}
	id := r.Form.Get("id")
	if len(id) == 0 {
		http.Error(w, "request parameters dont contain id", http.StatusBadRequest)
		return "", err
	}

	// Retrieve file
	kademliaID := model.KademliaIDFromString(id)
	file, err := lookupData(network, kademliaID)
	if err != nil {
		http.Error(w, "Error while retrieving file", http.StatusNotFound)
		return "", err
	}
	return string(file), nil
}

func store(w http.ResponseWriter, r *http.Request, network *model.KademliaNetwork) (string, error) {
	// Read input
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading request body", http.StatusBadRequest)
		return "", nil
	}
	// Store file
	file := []byte(string(body))
	id, err := storeData(network, file)
	if err != nil {
		http.Error(w, "error while storing file", http.StatusInternalServerError)
		return "", err
	}
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
