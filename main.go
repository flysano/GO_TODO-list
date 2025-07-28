package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("TODO_PORT") //получаем порт из переменной окружения
	if port == "" {
		port = "7540"
	}

	webDir := "./web" //определяем директорию для обслуживания файлов фронтенда

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	fmt.Println("Server started...")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}

}
