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

var topologyType = types.TopologyType_FiscoBcos
var attackDuration = 20

var FiscoBcosEvents = []*entities.Event{
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
		StartTime: time.Second * 150,
		Action:    types.ActionType_StartAttack,
		Handler:   func() error { return attack_manager.StartAttack(topologyType, attackDuration, 40) },
	},
	{
		StartTime: time.Second * 210,
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

func WithBlackListExperiment() error {
	configurationSettings := []*entities.ConfigurationSetting{
		{
			Mapping: map[string]string{},
		},
		{
			Mapping: map[string]string{},
		},
		{
			Mapping: map[string]string{},
		},
	}

	for _, configurationSetting := range configurationSettings {
		err := experiments.SingleSimulation(configurationSetting, FiscoBcosEvents)
		if err != nil {
			return fmt.Errorf("fiscobcos normal experiment failed: %v", err)
		}
	}

	return nil
}

/*
ttp server start successfully
{'attack_duration': 20, 'attack_thread_count': 40, 'attack_type': 'udp flood attack', 'attack_node': 'MaliciousNode-1', 'attacked_node': 'FabricOrderNode-3'}
10.134.86.192 - - [2025-10-28 10:46:41] "POST /startAttack HTTP/1.1" 200 129 0.008161
{'attack_duration': 20, 'attack_thread_count': 40, 'attack_type': 'udp flood attack', 'attack_node': 'MaliciousNode-1', 'attacked_node': 'FabricOrderNode-3'}
10.134.86.192 - - [2025-10-28 10:47:21] "POST /startAttack HTTP/1.1" 200 129 0.008350
*/

/*
{'attack_thread_count': 40, 'attack_type': 'udp flood attack', 'attack_node': 'MaliciousNode-1', 'attacked_node': 'FabricOrderNode-3', 'attack_duration': 20}
10.134.86.15 - - [2025-10-28 10:53:48] "POST /startAttack HTTP/1.1" 200 235 0.009020
*/
