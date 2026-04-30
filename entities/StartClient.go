package entities

// StartClient selected_network_layer: str, destination_port: int, processes: int, destinations: List[str], transmission_pattern: str
type StartClient struct {
	SelectedNetworkLayer string   `json:"selected_network_layer"`
	DestinationPort      int      `json:"destination_port"`
	Processes            int      `json:"processes"`
	Destinations         []string `json:"destinations"`
	TransmissionPattern  string   `json:"transmission_pattern"`
	FileSize             int      `json:"file_size"`
	BufferSize           int      `json:"buffer_size"`
	Content              string   `json:"content"` // 实际发送的内容

	// batch 之中需要指定的内容
	MessageSize int     `json:"message_size"`
	BatchSize   int     `json:"batch_size"`
	Interval    float64 `json:"interval"`
}

func NewStartClient(selectedNetworkLayer string, processes int, destinationPort int,
	destinations []string, transmissionPattern string, fileSize int,
	bufferSize int, content string, messageSize int,
	batchSize int, interval float64) *StartClient {
	return &StartClient{
		SelectedNetworkLayer: selectedNetworkLayer,
		DestinationPort:      destinationPort,
		Processes:            processes,
		Destinations:         destinations,
		TransmissionPattern:  transmissionPattern,
		FileSize:             fileSize,
		BufferSize:           bufferSize,
		Content:              content,

		MessageSize: messageSize,
		BatchSize:   batchSize,
		Interval:    interval,
	}
}
