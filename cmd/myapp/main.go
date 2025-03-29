package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"web-app/internal/web"
)

func main() {

	// Новый маршрутизатор
	mux := http.NewServeMux()

	// Обработка статических файлов
	staticDir := http.FileServer(http.Dir(filepath.Join("..", "..", "internal", "web", "templates")))

	mux.Handle("/generator.css", staticDir)
	mux.Handle("/success.html", staticDir)
	mux.Handle("/choose.css", staticDir)

	// Обработка маршрутов
	mux.HandleFunc("/", web.HomePageHandler)
	mux.HandleFunc("/send", web.GeneratePacketsHandler)
	mux.HandleFunc("/generator", web.Generator)

	// Запуск сервера
	port := ":8080"
	if len(os.Args) > 1 {
		port = ":" + os.Args[1]
	}
	log.Println("Starting server on", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
