package networking

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type messageAnswer struct {
	Message string `json:"message"`
}

// StartRestServer start the REST User API
func StartRestServer(exitChannel chan os.Signal) (*http.Server) {
	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/node/exit", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(
			messageAnswer{Message: "Server shutting down..."},
		)
		go func() {
			log.Print("Server shutting down...")
			exitChannel <- os.Interrupt
		}()
	})

	run := func() {
		log.Printf("RestServer ready ...")
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
	go run()
	return srv
}
