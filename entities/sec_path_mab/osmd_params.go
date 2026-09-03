package sec_path_mab

import "chain_simulation/entities/topology"

// SourceDestParams 图的源节点和目的节点的参数
type SourceDestParams struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

// SimNodeParam 模拟节点的参数
type SimNodeParam struct {
	Index                     int                        `json:"index"`
	Type                      string                     `json:"type"`
	CorruptRatio              topology.RatioDistribution `json:"corrupt_ratio"`
	CorruptSpecialPacketRatio topology.RatioDistribution `json:"corrupt_special_packet_ratio"`
}

type SimAbsLinkParam struct {
	SourceNode       SimNodeParam `json:"source_node"`       // 链路的源节点
	TargetNode       SimNodeParam `json:"target_node"`       // 链路的目的节点
	IntermediateNode SimNodeParam `json:"intermediate_node"` // pvlink 的中间节点
}

// SimLinkParam 模拟链路的参数
type SimLinkParam struct {
	SourceNode SimNodeParam `json:"source_node"` // 链路的源节点
	TargetNode SimNodeParam `json:"target_node"` // 链路的目的节点
}
