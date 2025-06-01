package models

type PacketParams struct {
	MacSrc       string
	MacDst       string
	IpSrc        string
	IpDst        string
	SrcPort      string
	DstPort      string
	TTL          string
	PacketSize   string
	LoopbackMode bool
}

type PacketInfo struct {
	Counter       uint64 `json:"counter"`       // Счетчик
	ForwardDelay  uint64 `json:"forwardDelay"`  // Задержка до шлейфа (SentTime -> RelayTime)
	BackwardDelay uint64 `json:"backwardDelay"` // Задержка обратно (RelayTime -> ReceivedTime)
	TotalDelay    uint64 `json:"totalDelay"`    // Общая задержка (SentTime -> ReceivedTime)
	InterArrival  uint64 `json:"interArrival"`  // Время между получением текущего и предыдущего пакета
}

type PacketReport struct {
	AverageTotal       float64  `json:"averageTotal"`        // Средняя общая задержка
	AverageForward     float64  `json:"averageForward"`      // Средняя задержка до шлейфа
	AverageBackward    float64  `json:"averageBackward"`     // Средняя задержка обратного пути
	MaxInterArrival    uint64   `json:"maxInterArrival"`     // Максимальная межпакетная задержка
	AverageInterArival float64  `json:"averageInterArrival"` // Средняя межпакетная задержка
	MissedPackets      []uint64 `json:"missedPackets"`       // Пропущенные пакеты
	LoopbackMode       bool     `json:"loopbackMode"`        // Флаг режима шлейфа
}

var LastReport PacketReport
