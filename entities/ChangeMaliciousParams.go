package entities

type ScheduledMaliciousParams struct {
	EmployedEpoch                  int `json:"employed_epoch"`
	NodeId                         int `json:"node_id"`
	CorruptRatioStart              int `json:"corrupt_ratio_start"`
	CorruptRatioEnd                int `json:"corrupt_ratio_end"`
	CorruptSpecialPacketRatioStart int `json:"corrupt_special_packet_ratio_start"`
	CorruptSpecialPacketRatioEnd   int `json:"corrupt_special_packet_ratio_end"`
}

func NewScheduledMaliciousParams(employedEpoch, nodeId, corruptRatioStart, corruptRatioEnd, corruptSpecialPacketRatioStart, corruptSpecialPacketRatioEnd int) *ScheduledMaliciousParams {
	return &ScheduledMaliciousParams{
		EmployedEpoch:                  employedEpoch,
		NodeId:                         nodeId,
		CorruptRatioStart:              corruptRatioStart,
		CorruptRatioEnd:                corruptRatioEnd,
		CorruptSpecialPacketRatioStart: corruptSpecialPacketRatioStart,
		CorruptSpecialPacketRatioEnd:   corruptSpecialPacketRatioEnd,
	}
}
