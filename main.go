package main

import (
	"GO_TODO-list/pkg/db"
	"GO_TODO-list/pkg/server"
	"log"
	"net/http"
	"os"
)

func main() {
	//подключение к БД
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	//определяем директорию для обслуживания файлов фронтенда
	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	//Запускаем сервер
	err := server.StartServer()
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
