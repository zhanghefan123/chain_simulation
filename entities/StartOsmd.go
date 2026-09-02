package entities

import "chain_simulation/entities/types"

type InitOsmd struct {
	// simulator params
	LearningRate         float64 `json:"learning_rate"`
	MinimumDeliveryRatio float64 `json:"minimum_delivery_ratio"`

	// client detailed info
	DestinationPort         int      `json:"destination_port"`
	Destinations            []string `json:"destinations"`
	MessageSize             int      `json:"message_size"`
	NumberOfPktsPerLink     int      `json:"number_of_pkts_per_link"`
	MiniBatchSize           int      `json:"mini_batch_size"`
	PacketSendingInterval   float64  `json:"packet_sending_interval"`
	SecPathMabStrategy      int      `json:"sec_path_mab_strategy"`
	EnableDadeAlgorithm     bool     `json:"enable_dade_algorithm"`
	EnableDedaAlgorithm     bool     `json:"enable_deda_algorithm"`
	MinAckForRttEstimation  int      `json:"min_ack_for_rtt_estimation"`
	ExperimentTimeElapsedMs int      `json:"experiment_time_elapsed_ms"`
}

func NewInitOsmd(learningRate float64, minimumDeliveryRatio float64,
	destinationPort int, destinations []string,
	messageSize int, numberOfPktsPerLink int,
	miniBatchSize int,
	packetSendingInterval float64, strategy types.SecPathMabStrategy,
	enableDadeAlgorithm bool, enableDedaAlgorithm bool,
	minAckForRttEstimation int, experimentTimeElapsedMs int) *InitOsmd {
	return &InitOsmd{
		LearningRate:         learningRate,
		MinimumDeliveryRatio: minimumDeliveryRatio,

		DestinationPort:         destinationPort,
		Destinations:            destinations,
		MessageSize:             messageSize,
		NumberOfPktsPerLink:     numberOfPktsPerLink,
		MiniBatchSize:           miniBatchSize,
		PacketSendingInterval:   packetSendingInterval,
		SecPathMabStrategy:      int(strategy),
		EnableDadeAlgorithm:     enableDadeAlgorithm,
		EnableDedaAlgorithm:     enableDedaAlgorithm,
		MinAckForRttEstimation:  minAckForRttEstimation,
		ExperimentTimeElapsedMs: experimentTimeElapsedMs,
	}
}
