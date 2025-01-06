package sender

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func SendPackets(interfaceName string, selected string, countOfPackets int, interval int, contentBytes []byte, identifiers []string) error {

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
		SrcMAC:       net.HardwareAddr(identifiers[0]),
		DstMAC:       net.HardwareAddr(identifiers[1]),
	}

	ip := layers.IPv4{
		Version:  4,
		TTL:      64,
		SrcIP:    net.IP(identifiers[2]),
		DstIP:    net.IP(identifiers[3]),
		Protocol: layers.IPProtocolUDP,
	}

	srcPort, err := strconv.Atoi(identifiers[4])
	if err != nil {
		fmt.Println("Ошибка преобразования в число порта источника")
		return err
	}

	dstPort, err := strconv.Atoi(identifiers[5])
	if err != nil {
		fmt.Println("Ошибка преобразования в число порта получателя")
		return err
	}

	udp := layers.UDP{
		SrcPort: layers.UDPPort(srcPort),
		DstPort: layers.UDPPort(dstPort),
	}

	fmt.Println(net.HardwareAddr(identifiers[0]))
	fmt.Println(net.HardwareAddr(identifiers[1]))
	fmt.Println(net.IP(identifiers[2]))
	fmt.Println(net.IP(identifiers[3]))
	fmt.Println(layers.UDPPort(srcPort))
	fmt.Println(layers.UDPPort(dstPort))

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
