package topology_manager

import (
	"fmt"

	"chain_simulation/configs"
	"chain_simulation/entities"
	"chain_simulation/entities/types"
	"chain_simulation/modules/internal/serviceapi"
)

var TopologyStartParamsInstance *entities.TopologyStartParams

func StartTopology(topologyType types.TopologyType, dynamicParameters *entities.DynamicParameters) error {
	params, err := entities.NewTopologyStartParams(topologyType, dynamicParameters)
	if err != nil {
		return fmt.Errorf("create topology parameters: %w", err)
	}
	if err = serviceapi.PostBackend(configs.TopConfigInstance.UrlConfig.StartTopologyUrl, params); err != nil {
		return fmt.Errorf("start topology: %w", err)
	}
	TopologyStartParamsInstance = params
	return nil
}

func StopTopology() error {
	if err := serviceapi.PostBackend(configs.TopConfigInstance.UrlConfig.StopTopologyUrl, nil); err != nil {
		return fmt.Errorf("stop topology: %w", err)
	}
	return nil
}
