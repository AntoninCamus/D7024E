package networking

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/LHJ/D7024E/kademlia"
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

	// Store
	if r.Method == "POST" {
		found = store(w, r)
	}

	// Find
	if r.Method == "GET" {
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

var results []string
func find(w http.ResponseWriter, r *http.Request) base64.Encoding {
	r.ParseForm()
	id := (r.Form.Get("id"))
	print(id)

	/* TODO lookup id and return file
	kademliaID =
	service.singleton.LookupData(kademliaID)
	*/

	var file base64.Encoding // Insert file content here. Call restService.singleton
	return file

}

func store(w http.ResponseWriter, r *http.Request) int {
	var results []string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
	}
	results = append(results, string(body))

	print(string(body))

	// TODO store file and return id


	var id int // insert id of file here
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
