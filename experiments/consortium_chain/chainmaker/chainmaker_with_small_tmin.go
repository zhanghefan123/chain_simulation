package chainmaker

import (
	"chain_simulation/entities"
	"chain_simulation/entities/types"
	"chain_simulation/experiments"
	"chain_simulation/modules/attack_manager"
	"chain_simulation/modules/consensus_manager"
	"chain_simulation/modules/topology_manager"
	"fmt"
	"time"
)

var ChainMakerWithSmallTminEvents = []*entities.Event{
	{
		StartTime: time.Second * 10,
		Action:    types.ActionType_StartTopology,
		Handler: func() error {
			return topology_manager.StartTopology(topologyType, &entities.DynamicParameters{ConsensusThreadCount: 20})
		},
	},
	{
		StartTime: time.Second * 90,
		Action:    types.ActionType_StartConsensus,
		Handler:   func() error { return consensus_manager.StartConsensus() },
	},
	{
		StartTime: time.Second * 110,
		Action:    types.ActionType_StartAttack,
		Handler:   func() error { return attack_manager.StartAttack(topologyType, attackDuration, 20) },
	},
	{
		StartTime: time.Second * 200,
		Action:    types.ActionType_StopConsensus,
		Handler:   func() error { return consensus_manager.StopConsensus() },
	},
	{
		StartTime: time.Second * 220,
		Action:    types.ActionType_StopTopology,
		Handler:   func() error { return topology_manager.StopTopology() },
	},
	{
		StartTime: time.Second * 240,
		Action:    types.ActionType_WaitTopologyRemove,
		Handler:   func() error { return nil },
	},
}

func ChainmakerWithSmallTminExperiment() error {
	configurationSettings := []*entities.ConfigurationSetting{
		{
			Mapping: map[string]string{
				"timeout_propose_optimal": "100",
				"propose_optimal":         "true",
				"enable_blacklist":        "false",
			},
		},
		{
			Mapping: map[string]string{
				"timeout_propose_optimal": "100",
				"propose_optimal":         "true",
				"enable_blacklist":        "false",
			},
		},
		{
			Mapping: map[string]string{
				"timeout_propose_optimal": "100",
				"propose_optimal":         "true",
				"enable_blacklist":        "false",
			},
		},
	}

	for _, configurationSetting := range configurationSettings {
		err := experiments.SingleSimulation(configurationSetting, ChainMakerWithSmallTminEvents)
		if err != nil {
			return fmt.Errorf("fisco bcos original experiment failed: %v", err)
		}
	}

	return nil
}
