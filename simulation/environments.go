package simulation

import (
	"chain_simulation/entities"
	"chain_simulation/entities/types"
	"chain_simulation/modules/attack_manager"
	"chain_simulation/modules/chaincode_manager"
	"chain_simulation/modules/consensus_manager"
	"chain_simulation/modules/topology_manager"
	"time"
)

var FabricEvents = []*entities.Event{
	{
		StartTime: time.Second * 10,
		Action:    types.ActionType_StartTopology,
		Handler:   func() error { return topology_manager.StartTopology(types.TopologyType_HyperledgerFabric) },
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
		Handler:   func() error { return attack_manager.StartAttack(types.TopologyType_HyperledgerFabric) },
	},
	{
		StartTime: time.Second * 200,
		Action:    types.ActionType_StopConsensus,
		Handler:   func() error { return consensus_manager.StopConsensus() },
	},
	{
		StartTime: time.Second * 205,
		Action:    types.ActionType_StopTopology,
		Handler:   func() error { return topology_manager.StopTopology() },
	},
	{
		StartTime: time.Second * 215,
		Action:    types.ActionType_WaitTopologyRemove,
		Handler:   func() error { return nil },
	},
}

var FiscoBcosEvents = []*entities.Event{
	{
		StartTime: time.Second * 10,
		Action:    types.ActionType_StartTopology,
		Handler:   func() error { return topology_manager.StartTopology(types.TopologyType_FiscoBcos) },
	},
	{
		StartTime: time.Second * 50,
		Action:    types.ActionType_StartConsensus,
		Handler:   func() error { return consensus_manager.StartConsensus() },
	},
	{
		StartTime: time.Second * 70,
		Action:    types.ActionType_StartAttack,
		Handler:   func() error { return attack_manager.StartAttack(types.TopologyType_FiscoBcos) },
	},
	{
		StartTime: time.Second * 150,
		Action:    types.ActionType_StopConsensus,
		Handler:   func() error { return consensus_manager.StopConsensus() },
	},
	{
		StartTime: time.Second * 155,
		Action:    types.ActionType_StopTopology,
		Handler:   func() error { return topology_manager.StopTopology() },
	},
	{
		StartTime: time.Second * 165,
		Action:    types.ActionType_WaitTopologyRemove,
		Handler:   func() error { return nil },
	},
}

var ChainMakerEvents = []*entities.Event{
	{
		StartTime: time.Second * 10,
		Action:    types.ActionType_StartTopology,
		Handler:   func() error { return topology_manager.StartTopology(types.TopologyType_ChainMaker) },
	},
	{
		StartTime: time.Second * 100,
		Action:    types.ActionType_StartConsensus,
		Handler:   func() error { return consensus_manager.StartConsensus() },
	},
	{
		StartTime: time.Second * 120,
		Action:    types.ActionType_StartAttack,
		Handler:   func() error { return attack_manager.StartAttack(types.TopologyType_ChainMaker) },
	},
	{
		StartTime: time.Second * 170,
		Action:    types.ActionType_StopConsensus,
		Handler:   func() error { return consensus_manager.StopConsensus() },
	},
	{
		StartTime: time.Second * 185,
		Action:    types.ActionType_StopTopology,
		Handler:   func() error { return topology_manager.StopTopology() },
	},
	{
		StartTime: time.Second * 205,
		Action:    types.ActionType_WaitTopologyRemove,
		Handler:   func() error { return nil },
	},
}
