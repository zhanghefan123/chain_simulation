package entities

import (
	"chain_simulation/configs"
	"chain_simulation/entities/types"
)

// 需要传递的参数
/*
let params = {
	attack_thread_count: selectedAttackThreadCount,
	attack_type: selectedAttackType, udp_flood_attack
	attack_node: selectedAttackNode.replace("_", "-"),
	attacked_node: selectedAttackedNode.replace("_", "-"),
	attack_duration: selectedAttackDuration
}
*/

type Attack struct {
	AttackDuration    int    `json:"attack_duration"`
	AttackThreadCount int    `json:"attack_thread_count"`
	AttackType        string `json:"attack_type"`
	AttackNode        string `json:"attack_node"`
	AttackedNode      string `json:"attacked_node"`
}

func NewAttackInstance(topologyType types.TopologyType) *Attack {
	attackInstance := &Attack{
		AttackDuration:    configs.TopConfigInstance.AttackConfig.Duration,
		AttackThreadCount: configs.TopConfigInstance.AttackConfig.ThreadCount,
		AttackType:        configs.TopConfigInstance.AttackConfig.Type,
		AttackNode:        "MaliciousNode-1",
	}
	switch topologyType {
	case types.TopologyType_HyperledgerFabric:
		attackInstance.AttackedNode = "FabricOrderNode-3"
	case types.TopologyType_FiscoBcos:
		attackInstance.AttackedNode = "FiscoBcosNode-3"
	case types.TopologyType_ChainMaker:
		attackInstance.AttackedNode = "ChainMakerNode-3"
	}
	return attackInstance
}
