package server

import (
	"GO_TODO-list/pkg/api"
	"fmt"
	"log"
	"net/http"
	"os"
)

func StartServer() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	api.Init()

	log.Println("Server started in port:", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return fmt.Errorf("Failed to start server: %w", err)
	}
	return nil
}
