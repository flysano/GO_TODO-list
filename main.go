package main

import (
	"GO_TODO-list/pkg/server"
	"fmt"
	"net/http"
)

func main() {
	webDir := "./web" //определяем директорию для обслуживания файлов фронтенда
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	//Запускаем сервер
	err := server.StartServer()
	if err != nil {
		fmt.Println(err)
	}

}
