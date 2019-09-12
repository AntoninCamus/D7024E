package networking

import (
	"encoding/json"
	"log"
	"net/http"
)

type messageAnswer struct {
	Message string `json:Message`
}

// StartRestServer start the REST User API
func StartRestServer() {
	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/node/exit", func(
		w http.ResponseWriter, r *http.Request,
	) {
		json.NewEncoder(w).Encode(
			messageAnswer{Message: "Server shutting down..."},
		)
		go func() {
			log.Fatal("Server shutting down...")
		}()
	})

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
