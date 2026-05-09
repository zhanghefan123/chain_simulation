package entities

import "chain_simulation/entities/types"

type SyncInstance struct {
	RateAdjustMode int `json:"rate_adjust_mode"`
}

func NewSyncInstance(rateAdjustMode types.RateAdjustMode) *SyncInstance {
	return &SyncInstance{
		RateAdjustMode: int(rateAdjustMode),
	}
}
