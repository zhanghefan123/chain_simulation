package fabrics

import (
	"chain_simulation/entities"
	"chain_simulation/entities/types"
	"chain_simulation/experiments"
	"chain_simulation/modules/attack_manager"
	"chain_simulation/modules/chaincode_manager"
	"chain_simulation/modules/consensus_manager"
	"chain_simulation/modules/topology_manager"
	"fmt"
	"time"
)

var topologyType = types.TopologyType_HyperledgerFabric

var FabricEvents = []*entities.Event{
	{
		StartTime: time.Second * 10,
		Action:    types.ActionType_StartTopology,
		Handler:   func() error { return topology_manager.StartTopology(topologyType) },
	},
	{
		StartTime: time.Second * 40,
		Action:    types.ActionType_StartInstallChaincode,
		Handler:   func() error { return chaincode_manager.InstallChainCode() },
	},
	{
		StartTime: time.Second * 80,
		Action:    types.ActionType_StartConsensus,
		Handler:   func() error { return consensus_manager.StartConsensus() },
	},
	{
		StartTime: time.Second * 100,
		Action:    types.ActionType_StartAttack,
		Handler:   func() error { return attack_manager.StartAttack(topologyType, 20) },
	},
	{
		StartTime: time.Second * 140,
		Action:    types.ActionType_StartAttack,
		Handler:   func() error { return attack_manager.StartAttack(topologyType, 20) },
	},
	{
		StartTime: time.Second * 200,
		Action:    types.ActionType_StopConsensus,
		Handler:   func() error { return consensus_manager.StopConsensus() },
	},
	{
		StartTime: time.Second * 210,
		Action:    types.ActionType_StopTopology,
		Handler:   func() error { return topology_manager.StopTopology() },
	},
	{
		StartTime: time.Second * 220,
		Action:    types.ActionType_WaitTopologyRemove,
		Handler:   func() error { return nil },
	},
}

func NormalExperiment() error {
	configurationSettings := []*entities.ConfigurationSetting{
		{
			Mapping: map[string]string{},
		},
	}

	for _, configurationSetting := range configurationSettings {
		err := experiments.SingleSimulation(configurationSetting, FabricEvents)
		if err != nil {
			return fmt.Errorf("fabric normal experiment failed: %v", err)
		}
	}
	return nil
}
