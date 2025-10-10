package consensus_manager

import (
	"chain_simulation/configs"
	"chain_simulation/utils/request"
	"fmt"
)

var ConsensusManagerInstance = &ConsensusManager{}

type ConsensusManager struct {
}

func StartConsensus() error {
	err := ConsensusManagerInstance.Start()
	if err != nil {
		return fmt.Errorf("start consensus manager failed: %v\n", err)
	}
	return nil
}

func StopConsensus() error {
	err := ConsensusManagerInstance.Stop()
	if err != nil {
		return fmt.Errorf("stop consensus manager failed: %v\n", err)
	}
	return nil
}

func (cm *ConsensusManager) Start() error {
	startConsensusUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		configs.TopConfigInstance.NetworkConfig.BackendPort,
		configs.TopConfigInstance.UrlConfig.StartTxRateTestUrl)
	err := request.PostJson(startConsensusUrl, nil)
	if err != nil {
		return fmt.Errorf("post topology failed: %v\n", err)
	}
	return nil
}

func (cm *ConsensusManager) Stop() error {
	stopConsensusUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		configs.TopConfigInstance.NetworkConfig.BackendPort,
		configs.TopConfigInstance.UrlConfig.StopTxRateTestUrl)
	err := request.PostJson(stopConsensusUrl, nil)
	if err != nil {
		return fmt.Errorf("post topology failed: %v\n", err)
	}
	return nil
}
