package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"web-api/internal/sender"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../internal/web/templates/home.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Could not execute template", http.StatusInternalServerError)
	}
}

func SendPacketsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	payloadSize, err := strconv.Atoi(r.FormValue("payloadSize"))
	if err != nil {
		fmt.Println("Ошибка преобразования:", err)
		http.Error(w, "Ошибка преобразования", http.StatusBadRequest)
		return
	}

	err = sender.SendPackets("eth0", payloadSize, 1_000_000_000_000) // Пример вызова функции
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("../../internal/web/templates/success.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка выполнени шабона", http.StatusInternalServerError)
	}
}
