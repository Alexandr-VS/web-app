package sender

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// ReceivePackets запускает приёмник пакетов и отправляет их в канал
func ReceivePackets(interfaceName string, packetChannel chan<- string, ipDst string, portDst string) {
	if handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever); err != nil {
		log.Fatalf("Ошибка открытия интерфейса: %v", err)
	} else if err = handle.SetBPFFilter("src host " + ipDst + " and udp port " + portDst); err != nil {
		log.Fatalf("Ошибка установки фильтра: %v", err)
	} else {
		defer handle.Close()

		packetSource := gopacket.NewPacketSource(handle, layers.LinkTypeEthernet)

		for packet := range packetSource.Packets() {
			// Обработка пакета
			if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
				ip, ok := ipLayer.(*layers.IPv4)
				if !ok {
					log.Println("Ошибка приведения типа для IP-слоя")
					continue
				}
				if packet.ApplicationLayer() != nil {
					// if ipDst == ip.SrcIP.String() {
					packetChannel <- fmt.Sprintf("IP-адрес источника: %s, IP-адрес получателя: %s, Полезная нагрузка: %d", ip.SrcIP, ip.DstIP, packet.Layer(layers.LayerTypeUDP).(*layers.UDP))
					// }
				}
			}
		}
	}
}
