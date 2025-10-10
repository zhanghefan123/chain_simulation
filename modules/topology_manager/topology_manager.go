package topology_manager

import (
	"chain_simulation/configs"
	"chain_simulation/entities"
	"chain_simulation/entities/types"
	"chain_simulation/utils/request"
	"fmt"
)

var TopologyManagerInstance = &TopologyManager{}

type TopologyManager struct {
}

func StartTopology(topologyType types.TopologyType) error {
	err := TopologyManagerInstance.Start(topologyType)
	if err != nil {
		return fmt.Errorf("start topology failed: %v", err)
	}
	return nil
}

func StopTopology() error {
	err := TopologyManagerInstance.Stop()
	if err != nil {
		return fmt.Errorf("stop topology failed: %v", err)
	}
	return nil
}

func (tm *TopologyManager) Start(topologyType types.TopologyType) error {
	startTopologyUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		configs.TopConfigInstance.NetworkConfig.BackendPort,
		configs.TopConfigInstance.UrlConfig.StartTopologyUrl)

	topology, err := entities.NewTopology(topologyType)
	if err != nil {
		return fmt.Errorf("create topology failed: %v\n", err)
	}
	fmt.Println(topology)
	err = request.PostJson(startTopologyUrl, topology)
	if err != nil {
		return fmt.Errorf("post topology failed: %v\n", err)
	}
	return nil
}

func (tm *TopologyManager) Stop() error {
	stopTopologyUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		configs.TopConfigInstance.NetworkConfig.BackendPort,
		configs.TopConfigInstance.UrlConfig.StopTopologyUrl)
	err := request.PostJson(stopTopologyUrl, nil)
	if err != nil {
		return fmt.Errorf("stop topology failed: %v\n", err)
	}
	return nil
}
