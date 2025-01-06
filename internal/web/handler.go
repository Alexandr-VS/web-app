package web

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"web-api/internal/sender"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	// Считывание шаблона
	tmpl, err := template.ParseFiles("../../internal/web/templates/home.html")
	if err != nil {
		log.Printf("Шаблон не найден: %v", err)
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

	macSrc := r.FormValue("mac-src")
	macDst := r.FormValue("mac-dst")
	ipSrc := r.FormValue("ip-src")
	ipDst := r.FormValue("ip-dst")
	srcPort := r.FormValue("src-port")
	dstPort := r.FormValue("dst-port")

	identifiers := []string{macSrc, macDst, ipSrc, ipDst, srcPort, dstPort}

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

	var contentBytes []byte
	if r.FormValue("dataSource") == "file" {
		err = r.ParseMultipartForm(10 << 20) // максимум 10Mb
		if err != nil {
			http.Error(w, "Ошибка при получении файла", http.StatusBadRequest)
			return
		}
		file, header, err := r.FormFile("filename")
		if err != nil {
			http.Error(w, "Ошибка при получении файла", http.StatusBadRequest)
			return
		}
		defer file.Close()

		if header.Size == 0 {
			http.Error(w, "Файл пуст", http.StatusBadRequest)
			return
		}

		contentBytes, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Ошибка при чтении файла", http.StatusInternalServerError)
			return
		}
	}

	err = sender.SendPackets("eth0", selected, countOfPackets, interval, contentBytes, identifiers)
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
