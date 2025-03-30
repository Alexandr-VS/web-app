package sender

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"time"
	"web-app/internal/models"
	"web-app/internal/utils"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func SendPackets(interfaceName string, selectedSrc string, countOfPackets int, interval float64, packetSizeStr string, contentBytes []byte, params models.PacketParams) error {
	// Открытие интерфейса для отправки пакетов
	handle, err := pcap.OpenLive(interfaceName, 1500, false, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("ошибка открытия интерфейса: %v", err)
	}
	defer handle.Close()

	buf := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	// Парсинг MAC и IP адресов
	srcMAC, err := utils.ParseMAC(params.MacSrc)
	if err != nil {
		return fmt.Errorf("ошибка парсинга MAC-адреса источника: %v", err)
	}

	dstMAC, err := utils.ParseMAC(params.MacDst)
	if err != nil {
		return fmt.Errorf("ошибка парсинга MAC-адреса получателя: %v", err)
	}

	srcIP, err := utils.ParseIP(params.IpSrc)
	if err != nil {
		return fmt.Errorf("ошибка парсинга IP-адреса источника: %v", err)
	}

	dstIP, err := utils.ParseIP(params.IpDst)
	if err != nil {
		return fmt.Errorf("ошибка парсинга IP-адреса получателя: %v", err)
	}

	// Преобразование портов и TTL
	srcPort, err := strconv.Atoi(params.SrcPort)
	if err != nil {
		return fmt.Errorf("ошибка преобразования порта источника: %v", err)
	}

	dstPort, err := strconv.Atoi(params.DstPort)
	if err != nil {
		return fmt.Errorf("ошибка преобразования порта получателя: %v", err)
	}

	ttl, err := strconv.Atoi(params.TTL)
	if err != nil {
		return fmt.Errorf("ошибка преобразования TTL: %v", err)
	}

	// Создание Ethernet, IP и UDP заголовков
	eth := layers.Ethernet{
		EthernetType: layers.EthernetTypeIPv4,
		SrcMAC:       srcMAC,
		DstMAC:       dstMAC,
	}

	ip := layers.IPv4{
		Version:  4,
		TTL:      uint8(ttl),
		SrcIP:    srcIP,
		DstIP:    dstIP,
		Protocol: layers.IPProtocolUDP,
	}

	udp := layers.UDP{
		SrcPort: layers.UDPPort(srcPort),
		DstPort: layers.UDPPort(dstPort),
	}

	udp.SetNetworkLayerForChecksum(&ip)

	// Отправка пакетов
	for i := 0; i < countOfPackets; i++ {
		payload, err := generatePayload(selectedSrc, packetSizeStr, contentBytes)
		if err != nil {
			return err
		}

		err = gopacket.SerializeLayers(buf, options,
			&eth,
			&ip,
			&udp,
			gopacket.Payload(payload),
		)
		if err != nil {
			return fmt.Errorf("ошибка сериализации пакетов: %v", err)
		}

		packetData := buf.Bytes()

		err = handle.WritePacketData(packetData)
		if err != nil {
			return fmt.Errorf("ошибка записи данных пакета: %v", err)
		}

		// Задержка между отправками пакетов
		time.Sleep(time.Duration(interval * float64(time.Second)))
	}
	return nil
}

// Функция для генерации полезной нагрузки
func generatePayload(selectedSrc string, packetSizeStr string, contentBytes []byte) ([]byte, error) {
	if selectedSrc == "pseudoRand" {
		var payloadSize *big.Int
		if packetSizeStr == "" {
			var err error
			payloadSize, err = rand.Int(rand.Reader, big.NewInt(1001)) // Генерация случайного размера до 1000
			if err != nil {
				return nil, fmt.Errorf("ошибка при генерации случайного размера пакета: %v", err)
			}
		} else {
			size, err := strconv.Atoi(packetSizeStr)
			if err != nil {
				return nil, fmt.Errorf("ошибка преобразования размера пакета: %v", err)
			}
			payloadSize = big.NewInt(int64(size))
		}

		payload := make([]byte, int(payloadSize.Int64()))
		_, err := rand.Read(payload) // заполнение полезной нагрузкой (вместо полинома происходит чтение из /dev/urandom или /dev/random, который вызывает getrandom(2))
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении случайных байтов: %v", err)
		}
		return payload, nil
	} else if selectedSrc == "file" {
		return contentBytes, nil
	}
	return nil, fmt.Errorf("неизвестный источник полезной нагрузки: %s", selectedSrc)
}
