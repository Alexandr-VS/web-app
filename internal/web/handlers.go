package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
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

var (
	packetChannel = make(chan models.PacketInfo)
	mu            sync.Mutex
	packets       []models.PacketInfo
)

func GetParamsToReceive(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../internal/web/templates/params.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка выполнения шаблона", http.StatusInternalServerError)
	}
}

func ReceivePacketsHandler(w http.ResponseWriter, r *http.Request) {
	var ipDst string
	var portDst string
	var totalPackets string
	if r.Method == http.MethodPost {
		mu.Lock()
		packets = []models.PacketInfo{}
		mu.Unlock()

		ipDst = r.FormValue("ip-dst")
		portDst = r.FormValue("port-dst")
		totalPackets = r.FormValue("totalPackets")

		// Запускаем приемник пакетов
		go sender.ReceivePackets("ens33", packetChannel, ipDst, portDst, totalPackets)
		go func() {
			for packet := range packetChannel {
				mu.Lock()
				packets = append(packets, packet) // добавление пакета в срез
				mu.Unlock()

				log.Printf("Получен пакет #%d отправленный в %d и принятый в %d", packet.Counter, packet.SentTime, packet.ReceivedTime)
			}
		}()

		tmpl, err := template.ParseFiles("../../internal/web/templates/receiver.html")

		if err != nil {
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)

		if err != nil {
			http.Error(w, "Ошибка выполнения шаблона", http.StatusInternalServerError)
		}
		return
	}

	// Обработка GET-запросов для получения списка пакетов
	if r.Method == http.MethodGet {
		mu.Lock()
		defer mu.Unlock()
		if err := json.NewEncoder(w).Encode(packets); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

func CheckCompletionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Проверяем, завершен ли прием пакетов
		completionStatus := map[string]bool{"completed": models.LastReport.AverageDelay > 0}
		w.Header().Set("Content-Type", "application/json") // Установка заголовка Content-Type
		if err := json.NewEncoder(w).Encode(completionStatus); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
}

func ReportHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("../../internal/web/templates/report.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Ошибка выполнения шаблона", http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodPost {
		// Отправляем последний отчет в формате JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(models.LastReport); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
