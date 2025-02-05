package sender

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"time"
	"web-app/internal/utils"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func SendPackets(interfaceName string, selectedSrc string, countOfPackets int, interval float64, contentBytes []byte, identifiers []string) error {

	handle, err := pcap.OpenLive(interfaceName, 1500, false, pcap.BlockForever)

	if err != nil {
		return err
	}
	defer handle.Close()

	buf := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	srcMAC, err := utils.ParseMAC(identifiers[0])
	if err != nil {
		return err
	}

	dstMAC, err := utils.ParseMAC(identifiers[1])
	if err != nil {
		return err
	}

	srcIP, err := utils.ParseIP(identifiers[2])
	if err != nil {
		return err
	}

	dstIP, err := utils.ParseIP(identifiers[3])
	if err != nil {
		return err
	}

	srcPort, err := strconv.Atoi(identifiers[4])
	if err != nil {
		return fmt.Errorf("ошибка преобразования в число порта источника: %v", err)
	}

	dstPort, err := strconv.Atoi(identifiers[5])
	if err != nil {
		return fmt.Errorf("ошибка преобразования в число порта получателя: %v", err)
	}

	ttl, err := strconv.Atoi(identifiers[6])
	if err != nil {
		return fmt.Errorf("ошибка преобразовани в число ttl")
	}

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

	var payload []byte

	for i := 0; i < countOfPackets; i++ {
		if selectedSrc == "pseudoRand" {
			payloadSize, err := rand.Int(rand.Reader, big.NewInt(1001))
			if err != nil {
				fmt.Println("Ошибка при генерации случайного числа:", err)
				return err
			}
			payload = make([]byte, int(payloadSize.Int64()))
			_, err = rand.Read(payload)
			if err != nil {
				fmt.Println("Ошибка при чтении случайных байтов:", err)
				return err
			}
		} else if selectedSrc == "file" {
			payload = contentBytes
		}

		err = gopacket.SerializeLayers(buf, options,
			&eth,
			&ip,
			&udp,
			gopacket.Payload(payload),
		)

		if err != nil {
			return err
		}

		packetData := buf.Bytes()

		err = handle.WritePacketData(packetData)

		if err != nil {
			return err
		}

		time.Sleep(time.Duration(interval * float64(time.Second)))
	}
	return nil
}
