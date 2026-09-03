package consensus_manager

import (
	"fmt"

	"chain_simulation/configs"
	"chain_simulation/modules/internal/serviceapi"
)

func StartConsensus() error {
	if err := serviceapi.PostBackend(configs.TopConfigInstance.UrlConfig.StartTxRateTestUrl, nil); err != nil {
		return fmt.Errorf("start consensus: %w", err)
	}
	return nil
}

func StopConsensus() error {
	if err := serviceapi.PostBackend(configs.TopConfigInstance.UrlConfig.StopTxRateTestUrl, nil); err != nil {
		return fmt.Errorf("stop consensus: %w", err)
	}
	return nil
}
