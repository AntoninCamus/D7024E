package networking

import (
	"encoding/json"
	"github.com/LHJ/D7024E/kademlia"
	"log"
	"net/http"
	"fmt"
)

type messageAnswer struct {
	Message string `json:"message"`
}

type restService struct {
	server *http.Server
	singleton *kademlia.Kademlia
}

// StartRestServer start the REST User API
func StartRestServer(s *kademlia.Kademlia) (*restService) {
	fmt.Println("Starting server...")

	service := restService{
		server:    &http.Server{Addr: ":8080", Handler: nil},
		singleton: s,
	}
	http.HandleFunc("/kademlia/store", service.store)
	http.HandleFunc("/kademlia/find", service.find)
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

func (s *restService) store(w http.ResponseWriter, r *http.Request) {

}
type test_struct struct {
	Test string
}


func (s *restService) find(w http.ResponseWriter, r *http.Request) {

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
