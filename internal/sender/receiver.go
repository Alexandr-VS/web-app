package sender

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"
	"time"
	"web-app/internal/models"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// ReceivePackets запускает приёмник пакетов и отправляет их в канал
func ReceivePackets(interfaceName string, packetChannel chan<- models.PacketInfo, ipDst string, portDst string, totalPacketsStr string, loopbackMode bool) {
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

	var (
		totalForwardDelay  uint64
		totalBackwardDelay uint64
		maxInterArrival    uint64
		lastReceivedTime   uint64
		receivedPackets    uint64
		totalDelay         uint64
		totalInterArrival  uint64
	)
	missedPackets := make([]uint64, 0)

	receivedMap := make(map[uint64]bool)

	// Таймер для остановки приёмника
	timer := time.NewTimer(5 * time.Second)
	if loopbackMode {
		timer = time.NewTimer(10 * time.Second)
	}
	defer timer.Stop()

	for {
		select {
		case packet := <-packetSource.Packets():
			// Сброс таймера при получении пакета
			if !timer.Stop() {
				<-timer.C // Прочитать, если таймер уже сработал
			}
			timer.Reset(7 * time.Second)

			// Обработка пакета

			if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
				_, ok := ipLayer.(*layers.IPv4)
				if !ok {
					log.Println("Ошибка приведения типа для IP-слоя")
					continue
				}

				if appLayer := packet.ApplicationLayer(); appLayer != nil {
					payload := appLayer.Payload()
					if len(payload) < 24 {
						log.Printf("Пакет слишком мал: %d байт\n", len(payload))
						continue
					}

					counter := binary.BigEndian.Uint64(payload[0:8])
					sentTime1 := binary.BigEndian.Uint64(payload[8:16])
					sentTime2 := binary.BigEndian.Uint64(payload[16:24]) // Время отправки от шлейфа
					receivedTime := uint64(time.Now().UnixNano())

					packetInfo := models.PacketInfo{
						Counter:    counter,
						TotalDelay: (receivedTime - sentTime1) / 1e6,
					}

					packetInfo.InterArrival = 0
					if lastReceivedTime > 0 {
						packetInfo.InterArrival = receivedTime - lastReceivedTime
					}

					if packetInfo.InterArrival > maxInterArrival {
						maxInterArrival = packetInfo.InterArrival
					}

					lastReceivedTime = receivedTime
					if loopbackMode {
						forwardDelayNs := int64(sentTime2) - int64(sentTime1)
						backwardDelayNs := int64(receivedTime) - int64(sentTime2)

						packetInfo.ForwardDelay = uint64(math.Abs(float64(forwardDelayNs) / 1e6))
						packetInfo.BackwardDelay = uint64(math.Abs(float64(backwardDelayNs) / 1e6))
					}

					packetChannel <- packetInfo

					// Накопление сумм для среднего
					totalForwardDelay += packetInfo.ForwardDelay
					totalBackwardDelay += packetInfo.BackwardDelay
					totalDelay += packetInfo.TotalDelay
					totalInterArrival += packetInfo.InterArrival

					receivedPackets++

					// Добавляем номер пакета в хеш-таблицу
					receivedMap[counter] = true

					fmt.Println("Пакет номер:", packetInfo.Counter,
						"Время отправки с генератора:", sentTime1,
						"Время отправки с шлейфа:", sentTime2,
						"Время приёма пакета:", receivedTime,
						"До шлейфа:", packetInfo.ForwardDelay,
						"От шлейфа:", packetInfo.BackwardDelay)
				}
			}
		case <-timer.C:
			log.Println("Таймер сработал, прекращаем приём пакетов")
			close(packetChannel)
			// Вычисление среднего времени задержки
			var averageDelay float64
			var averageForward float64
			var averageBackward float64
			var averageInterArrival float64

			if receivedPackets > 0 {
				averageDelay = float64(totalDelay) / float64(receivedPackets)
				averageForward = float64(totalForwardDelay) / float64(receivedPackets)
				averageBackward = float64(totalBackwardDelay) / float64(receivedPackets)
				averageInterArrival = float64(totalInterArrival) / float64(receivedPackets)
			}

			totalPackets, err := strconv.Atoi(totalPacketsStr)
			if err != nil {
				log.Println("Ошибка преобразования количества пакетов:", err)
			}

			// Определение пропущенных пакетов
			for i := uint64(0); i < uint64(totalPackets); i++ {
				if !receivedMap[i] {
					missedPackets = append(missedPackets, i)
				}
			}

			models.LastReport = models.PacketReport{
				AverageTotal:       averageDelay,
				MissedPackets:      missedPackets,
				LoopbackMode:       loopbackMode,
				AverageForward:     averageForward,
				AverageBackward:    averageBackward,
				MaxInterArrival:    maxInterArrival,
				AverageInterArival: averageInterArrival,
			}
		}
	}
}

var (
	relayStopChan chan struct{}
	relayMu       sync.Mutex
)

func StartRelay(interfaceName, ipDst, portDst string) {
	relayMu.Lock()
	defer relayMu.Unlock()

	// Если уже запущен — остановим старый
	if relayStopChan != nil {
		close(relayStopChan)
		relayStopChan = nil
	}

	relayStopChan = make(chan struct{})

	go ReceivePacketsWithRelay(interfaceName, ipDst, portDst, relayStopChan)
}

func StopRelay() {
	relayMu.Lock()
	defer relayMu.Unlock()

	if relayStopChan != nil {
		close(relayStopChan)
		relayStopChan = nil
	}
}

// ReceivePacketsWithRelay принимает пакеты, меняет IP, порты и MAC, и отправляет обратно
func ReceivePacketsWithRelay(interfaceName, ipDst, portDst string, stopChan chan struct{}) {
	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatalf("Ошибка открытия интерфейса: %v", err)
	}
	defer handle.Close()
	log.Println("Приём с ретрансляцией начат")

	filter := fmt.Sprintf("dst host %s and udp dst port %s", ipDst, portDst)
	if err := handle.SetBPFFilter(filter); err != nil {
		log.Fatalf("Ошибка установки фильтра: %v", err)
	}

	packetSource := gopacket.NewPacketSource(handle, layers.LinkTypeEthernet)

	for {
		select {
		case <-stopChan:
			log.Println("Режим шлейфа остановлен")
			return
		case packet, ok := <-packetSource.Packets():
			if !ok {
				log.Println("Источник пакетов закрыт")
				return
			}
			// Обработка пакета (ваш код)
			ethLayer := packet.Layer(layers.LayerTypeEthernet)
			ipLayer := packet.Layer(layers.LayerTypeIPv4)
			udpLayer := packet.Layer(layers.LayerTypeUDP)

			if ethLayer == nil || ipLayer == nil || udpLayer == nil {
				continue
			}

			eth, _ := ethLayer.(*layers.Ethernet)
			ip, _ := ipLayer.(*layers.IPv4)
			udp, _ := udpLayer.(*layers.UDP)

			ip.SrcIP, ip.DstIP = ip.DstIP, ip.SrcIP
			udp.SrcPort, udp.DstPort = udp.DstPort, udp.SrcPort
			eth.SrcMAC, eth.DstMAC = eth.DstMAC, eth.SrcMAC

			buf := gopacket.NewSerializeBuffer()
			opts := gopacket.SerializeOptions{
				ComputeChecksums: true,
				FixLengths:       true,
			}

			timestamp := uint64(time.Now().UnixNano())
			binary.BigEndian.PutUint64(udp.Payload[16:24], timestamp)

			udp.SetNetworkLayerForChecksum(ip)

			err := gopacket.SerializeLayers(buf, opts,
				eth,
				ip,
				udp,
				gopacket.Payload(udp.Payload),
			)
			if err != nil {
				log.Println("Ошибка сериализации пакета:", err)
				continue
			}

			err = handle.WritePacketData(buf.Bytes())
			if err != nil {
				log.Println("Ошибка отправки пакета:", err)
			}
		}
	}
}
