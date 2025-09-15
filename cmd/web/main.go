package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

var handler http.HandlerFunc

func main() {
	router := chi.NewRouter()

	handler = func(writer http.ResponseWriter, request *http.Request) {
		hostname, err := os.Hostname()
		if err != nil {
			fmt.Printf("Error getting hostname: %v\n", err)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(map[string]interface{}{
			"hostname": hostname,
			"time":     time.Now(),
			"message":  "New attribute message",
		})
	}

	router.Method(http.MethodGet, "/api/v1/details", handler)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Printf("Starting subscriptions service on port %s", server.Addr)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}

}
