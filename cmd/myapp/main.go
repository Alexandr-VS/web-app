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

	mux.Handle("/choose.css", staticDir)
	mux.Handle("/generator.css", staticDir)
	mux.Handle("/receiver.css", staticDir)
	mux.Handle("/generator.js", staticDir)
	mux.Handle("/receiver.js", staticDir)
	mux.Handle("/params.js", staticDir)

	// Обработка маршрутов
	mux.HandleFunc("/", web.HomePageHandler)
	mux.HandleFunc("/send", web.GeneratePacketsHandler)
	mux.HandleFunc("/generator", web.Generator)
	mux.HandleFunc("/receiver", web.ReceivePacketsHandler)
	mux.HandleFunc("/params", web.GetParamsToReceive)
	mux.HandleFunc("/report", web.ReportHandler)
	mux.HandleFunc("/check-completion", web.CheckCompletionHandler)

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
