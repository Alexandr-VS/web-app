package sender

import (
	"crypto/rand"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func SendPackets(interfaceName string, payloadSize int, desiredBandwidth int64) error {
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

	payload := make([]byte, payloadSize)
	rand.Read(payload)

	err = gopacket.SerializeLayers(buf, options,
		&eth,
		&ip,
		&udp,
		gopacket.Payload(payload),
	)

	if err != nil {
		return err
	}

	packetSizeInBits := int64(payloadSize * 8)
	packetsPerSecond := desiredBandwidth / packetSizeInBits
	interval := time.Second / time.Duration(packetsPerSecond)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	packetData := buf.Bytes()

	// for range ticker.C {
	// 	err = handle.WritePacketData(packetData)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	err = handle.WritePacketData(packetData)
	if err != nil {
		return err
	}

	return nil
}
