package models

type PacketParams struct {
	MacSrc     string
	MacDst     string
	IpSrc      string
	IpDst      string
	SrcPort    string
	DstPort    string
	TTL        string
	PacketSize string
}

type PacketInfo struct {
	Counter      uint64 `json:"counter"`
	SentTime     uint64 `json:"sentTime"`
	ReceivedTime uint64 `json:"receivedTime"`
	Delay        uint64 `json:"delay"`
}

type PacketReport struct {
	AverageDelay  float64  `json:"averageDelay"`
	MissedPackets []uint64 `json:"missedPackets"`
}

var LastReport PacketReport
