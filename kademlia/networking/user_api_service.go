package networking

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/LHJ/D7024E/kademlia"
	"github.com/LHJ/D7024E/kademlia/model"
	"io/ioutil"
	"log"
	"net/http"
)

type messageAnswer struct {
	Message string `json:"message"`
}

type restService struct {
	server *http.Server
	singleton *kademlia.Kademlia
}
var service restService
// StartRestServer start the REST User API
func StartRestServer(s *kademlia.Kademlia) (*restService) {
	fmt.Println("Starting server...")

	service := restService{
		server:    &http.Server{Addr: ":8080", Handler: nil},
		singleton: s,
	}
	http.HandleFunc("/kademlia/file", service.findstore)
	http.HandleFunc("/node/exit", service.exitServer)

	serving := func() {
		err := service.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
	go serving()

	return &service
}

func (s *restService) findstore(w http.ResponseWriter, r *http.Request) {
	var found interface{} // Generic

	if r.Method == "POST" { // Store
		found = store(w, r)
	} else if r.Method == "GET" { // Find
		found = find(w, r)
	}

	// Write response
	js, err := json.Marshal(found)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func find(w http.ResponseWriter, r *http.Request) base64.Encoding {
	// Read input
	r.ParseForm()
	id := (r.Form.Get("id"))
	print(id)

	// Retrieve file
	kademliaID := model.NewKademliaID(id)
	file := service.singleton.LookupData(kademliaID)
	return file

}

func store(w http.ResponseWriter, r *http.Request) string {
	// Read input
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
	}

	// Store file
	file := []byte(string(body))
	id, err := service.singleton.Store(file)
	return id
}


func (s *restService) exitServer(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(
		messageAnswer{Message: "Server shutting down..."},
	)

	go func() {
		_ = s.server.Close()
		log.Fatal("Server shutting down...")
	}()
}
