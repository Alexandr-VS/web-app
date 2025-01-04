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
	// Считывание шаблона
	tmpl, err := template.ParseFiles("../../internal/web/templates/home.html")
	if err != nil {
		log.Printf("Шаблоне не найден: %v", err)
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	// Запуск шаблона
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка выполнения шабона", http.StatusInternalServerError)
	}
}

func SendPacketsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	selected := r.FormValue("dataSource")

	countOfPackets, err := strconv.Atoi(r.FormValue("countOfPackets"))
	if err != nil {
		fmt.Println("Ошибка преобразования количества пакетов:", err)
		http.Error(w, "Ошибка преобразования количества пакетов", http.StatusBadRequest)
		return
	}

	interval, err := strconv.Atoi(r.FormValue("interval"))
	if err != nil {
		fmt.Println("Ошибка преобразования:", err)
		http.Error(w, "Ошибка преобразования", http.StatusBadRequest)
		return
	}

	err = sender.SendPackets("eth0", selected, countOfPackets, interval)
	if err != nil {
		fmt.Println("Ошибка отправки")
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
		http.Error(w, "Ошибка выполнения шаблона", http.StatusInternalServerError)
	}
}
