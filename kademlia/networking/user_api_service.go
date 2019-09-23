package networking

import (
	"encoding/json"
	"fmt"
	"github.com/LHJ/D7024E/kademlia"
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
	found := ""

	if r.Method == "POST" {
		found = "insert id of file here"
	}
	if r.Method == "GET" {
		found = "insert content of file here"
	}

	js, err := json.Marshal(found)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)


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
