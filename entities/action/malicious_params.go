package action

type ScheduledMaliciousParams struct {
	EmployedEpochOrTimestampMs     int `json:"employed_epoch_or_timestamp_ms"`
	NodeId                         int `json:"node_id"`
	CorruptRatioStart              int `json:"corrupt_ratio_start"`
	CorruptRatioEnd                int `json:"corrupt_ratio_end"`
	CorruptSpecialPacketRatioStart int `json:"corrupt_special_packet_ratio_start"`
	CorruptSpecialPacketRatioEnd   int `json:"corrupt_special_packet_ratio_end"`
}

func NewScheduledMaliciousParams(employedEpochOrTimestampMs, nodeId, corruptRatioStart, corruptRatioEnd, corruptSpecialPacketRatioStart, corruptSpecialPacketRatioEnd int) *ScheduledMaliciousParams {
	return &ScheduledMaliciousParams{
		EmployedEpochOrTimestampMs:     employedEpochOrTimestampMs,
		NodeId:                         nodeId,
		CorruptRatioStart:              corruptRatioStart,
		CorruptRatioEnd:                corruptRatioEnd,
		CorruptSpecialPacketRatioStart: corruptSpecialPacketRatioStart,
		CorruptSpecialPacketRatioEnd:   corruptSpecialPacketRatioEnd,
	}
}
