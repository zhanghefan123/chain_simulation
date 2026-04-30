package fiscobcos

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

var FiscoBcosOriginalEvents = []*entities.Event{
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
		Handler:   func() error { return attack_manager.StartAttack(topologyType, attackDuration, 40) },
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

func OriginalExperiment() error {
	configurationSettings := []*entities.ConfigurationSetting{
		{
			Mapping: map[string]string{
				"fast_trigger":                "true",
				"recursive_trigger":           "true",
				"use_modified_request_blocks": "false",
				"enable_pending":              "false",
				"enable_blacklist":            "false",
			},
		},
		{
			Mapping: map[string]string{
				"fast_trigger":                "true",
				"recursive_trigger":           "true",
				"use_modified_request_blocks": "false",
				"enable_pending":              "false",
				"enable_blacklist":            "false",
			},
		},
		{
			Mapping: map[string]string{
				"fast_trigger":                "true",
				"recursive_trigger":           "true",
				"use_modified_request_blocks": "false",
				"enable_pending":              "false",
				"enable_blacklist":            "false",
			},
		},
		{
			Mapping: map[string]string{
				"fast_trigger":                "true",
				"recursive_trigger":           "true",
				"use_modified_request_blocks": "false",
				"enable_pending":              "false",
				"enable_blacklist":            "false",
			},
		},
		{
			Mapping: map[string]string{
				"fast_trigger":                "true",
				"recursive_trigger":           "true",
				"use_modified_request_blocks": "false",
				"enable_pending":              "false",
				"enable_blacklist":            "false",
			},
		},
		{
			Mapping: map[string]string{
				"fast_trigger":                "true",
				"recursive_trigger":           "true",
				"use_modified_request_blocks": "false",
				"enable_pending":              "false",
				"enable_blacklist":            "false",
			},
		},
	}

	for _, configurationSetting := range configurationSettings {
		err := experiments.SingleSimulation(configurationSetting, FiscoBcosOriginalEvents)
		if err != nil {
			return fmt.Errorf("fisco bcos original experiment failed: %v", err)
		}
	}

	return nil
}
