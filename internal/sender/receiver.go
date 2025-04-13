package sender

import (
	"encoding/binary"
	"log"
	"strconv"
	"time"
	"web-app/internal/models"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// ReceivePackets запускает приёмник пакетов и отправляет их в канал
func ReceivePackets(interfaceName string, packetChannel chan<- models.PacketInfo, ipDst string, portDst string, totalPacketsStr string) {
	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatalf("Ошибка открытия интерфейса: %v", err)
	}
	defer handle.Close()
	log.Println("Приём начат")

	if err := handle.SetBPFFilter("src host " + ipDst + " and udp port " + portDst); err != nil {
		log.Fatalf("Ошибка установки фильтра: %v", err)
	}

	packetSource := gopacket.NewPacketSource(handle, layers.LinkTypeEthernet)

	var totalDelay uint64
	var receivedPackets uint64
	missedPackets := make([]uint64, 0)

	receivedMap := make(map[uint64]bool)

	// Таймер для остановки приёмника
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	for {
		select {
		case packet := <-packetSource.Packets():
			// Сброс таймера при получении пакета
			if !timer.Stop() {
				<-timer.C // Прочитать, если таймер уже сработал
			}
			timer.Reset(5 * time.Second)

			// Обработка пакета
			receivedTime := uint64(time.Now().UnixMilli())

			if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
				_, ok := ipLayer.(*layers.IPv4)
				if !ok {
					log.Println("Ошибка приведения типа для IP-слоя")
					continue
				}

				if appLayer := packet.ApplicationLayer(); appLayer != nil {
					payload := appLayer.Payload()
					if len(payload) < 16 {
						log.Printf("Пакет слишком мал: %d байт\n", len(payload))
						continue
					}

					counter := binary.BigEndian.Uint64(payload[0:8])
					sentTime := binary.BigEndian.Uint64(payload[8:16])

					packetChannel <- models.PacketInfo{
						Counter:      counter,
						SentTime:     sentTime,
						ReceivedTime: receivedTime,
						Delay:        receivedTime - sentTime,
					}
					totalDelay += receivedTime - sentTime
					receivedPackets++

					// Добавляем номер пакета в хеш-таблицу
					receivedMap[counter] = true
				}
			}
		case <-timer.C:
			log.Println("Таймер сработал, прекращаем приём пакетов")
			close(packetChannel)
			// Вычисление среднего времени задержки
			var averageDelay float64
			if receivedPackets > 0 {
				averageDelay = float64(totalDelay) / float64(receivedPackets)
			}

			totalPackets, err := strconv.Atoi(totalPacketsStr)
			if err != nil {
				log.Println("Ошибка преобразования количества пакетов:", err)
			}

			// Определение пропущенных пакетов
			for i := uint64(0); i < uint64(totalPackets); i++ {
				if i >= receivedPackets {
					missedPackets = append(missedPackets, i)
				}
			}

			models.LastReport = models.PacketReport{
				AverageDelay:  averageDelay,
				MissedPackets: missedPackets,
			}
		}
	}
}
