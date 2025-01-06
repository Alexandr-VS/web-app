package main

import (
	"log"
	"net/http"
	"web-api/internal/web"
)

func main() {

	// Новый маршрутизатор
	mux := http.NewServeMux()

	// Обработка маршрутов
	mux.HandleFunc("/", web.HomePageHandler)
	mux.HandleFunc("/send", web.SendPacketsHandler)

	// Запуск сервера
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
