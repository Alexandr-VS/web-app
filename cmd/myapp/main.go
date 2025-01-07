package main

import (
	"log"
	"net/http"
	"path/filepath"
	"web-api/internal/web"
)

func main() {

	// Новый маршрутизатор
	mux := http.NewServeMux()

	// Обработка статических файлов
	staticDir := http.FileServer(http.Dir(filepath.Join("..", "..", "internal", "web", "templates")))

	mux.Handle("/styles.css", staticDir)   // Обработка CSS файла
	mux.Handle("/success.html", staticDir) // Обработка success.html, если нужно

	// Обработка маршрутов
	mux.HandleFunc("/", web.HomePageHandler)
	mux.HandleFunc("/send", web.SendPacketsHandler)

	// Запуск сервера
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
