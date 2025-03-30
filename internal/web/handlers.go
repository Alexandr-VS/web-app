package web

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"web-app/internal/models"
	"web-app/internal/sender"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	// Считывание шаблона
	tmpl, err := template.ParseFiles("../../internal/web/templates/choose.html")
	if err != nil {
		log.Printf("Шаблон не найден: %v", err)
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	// Запуск шаблона
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка выполнения шаблона", http.StatusInternalServerError)
	}
}

func Generator(w http.ResponseWriter, r *http.Request) {
	// Считывание шаблона
	tmpl, err := template.ParseFiles("../../internal/web/templates/generator.html")
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

var packetChannel = make(chan string)

// Обработчик отправки пакетов
func GeneratePacketsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	params := models.PacketParams{
		MacSrc:     r.FormValue("mac-src"),
		MacDst:     r.FormValue("mac-dst"),
		IpSrc:      r.FormValue("ip-src"),
		IpDst:      r.FormValue("ip-dst"),
		SrcPort:    r.FormValue("src-port"),
		DstPort:    r.FormValue("dst-port"),
		TTL:        r.FormValue("TTL"),
		PacketSize: r.FormValue("packetSize"),
	}

	if params.TTL == "" {
		params.TTL = "64"
	}

	selectedSrc := r.FormValue("dataSource")

	countOfPackets, err := strconv.Atoi(r.FormValue("countOfPackets"))
	if err != nil {
		fmt.Println("Ошибка преобразования количества пакетов:", err)
		http.Error(w, "Ошибка преобразования количества пакетов", http.StatusBadRequest)
		return
	}

	interval, err := strconv.ParseFloat(r.FormValue("interval"), 64)
	if err != nil {
		fmt.Println("Ошибка преобразования интервала:", err)
		http.Error(w, "Ошибка преобразования интервала", http.StatusBadRequest)
		return
	}

	packetSizeStr := r.FormValue("packetSize")

	var contentBytes []byte

	if selectedSrc == "file" {
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

	if r.FormValue("toggleSwitch") == "on" {
		// режим шлейфа
	}

	go sender.SendPackets("ens33", selectedSrc, countOfPackets, interval, packetSizeStr, contentBytes, params)
	if err != nil {
		fmt.Println("Ошибка отправки пакетов")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/generator", http.StatusSeeOther)
}

func ReceivePacketsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../internal/web/templates/receiver.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	// Отправляем данные из канала в шаблон
	go func() {
		for packet := range packetChannel {
			// Здесь можно обновить состояние для отображения на странице
			log.Println(packet) // Логируем полученные пакеты
			// Вы можете использовать механизм обновления страницы через AJAX или WebSocket
		}
	}()

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка выполнения шаблона", http.StatusInternalServerError)
	}
}
