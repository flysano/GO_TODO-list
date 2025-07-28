package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func StartServer() error {
	port := os.Getenv("TODO_PORT") //получаем порт из переменной окружения
	if port == "" {
		port = "7540"
	}

	log.Println("Server started in port:", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return fmt.Errorf("Failed to start server: %w", err)
	}
	return nil
}
