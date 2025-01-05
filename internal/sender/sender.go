package sender

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func SendPackets(interfaceName string, selected string, countOfPackets int, interval int, contentBytes []byte) error {

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

	eth := layers.Ethernet{
		EthernetType: layers.EthernetTypeIPv4,
		SrcMAC:       net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		DstMAC:       net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}

	ip := layers.IPv4{
		Version:  4,
		TTL:      64,
		SrcIP:    net.IP{127, 0, 0, 1},
		DstIP:    net.IP{127, 0, 0, 1},
		Protocol: layers.IPProtocolUDP,
	}

	udp := layers.UDP{
		SrcPort: 62003,
		DstPort: 8080,
	}

	udp.SetNetworkLayerForChecksum(&ip)

	var payload []byte

	for i := 0; i < countOfPackets; i++ {
		if selected == "pseudoRand" {
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
		} else if selected == "file" {
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
		time.Sleep(time.Duration(interval * int(time.Second)))
	}
	return nil
}
