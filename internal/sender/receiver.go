package sender

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// ReceivePackets запускает приёмник пакетов и отправляет их в канал
func ReceivePackets(interfaceName string, packetChannel chan<- string) {
	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatalf("Ошибка открытия интерфейса: %v", err)
	}
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
			packetChannel <- fmt.Sprintf("IP-адрес отправителя: %s, IP-адрес получателя: %s", ip.SrcIP, ip.DstIP)
		} else {
			log.Println("Пакет не содержит IP-слоя")
		}
	}
}
