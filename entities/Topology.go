package entities

import (
	"chain_simulation/entities/types"
	"chain_simulation/utils/file"
	"encoding/json"
	"fmt"
)

var (
	TopologyPathFabric              = "../resources/topologies/fabric_test_topology.json"
	TopologyPathFisco               = "../resources/topologies/fisco_bcos_test_topology.json"
	TopologyPathChainmaker          = "../resources/topologies/chainmaker_test_topology.json"
	TopologySimplePathValidation    = "../resources/topologies/simple_path_validation_topology.json"
	TopologyMulticastPathValidation = "../resources/topologies/multicast_path_validation_topology.json"
	TopologySecPathMab              = "../resources/topologies/sec_path_mab_topology.json"
)

type Topology struct {
	ConsensusThreadCount int `json:"consensus_thread_count"`
	AccessLinkBandwidth  int `json:"access_link_bandwidth"`
	ConsensusNodeCpu     int `json:"consensus_node_cpu"`
	ConsensusNodeMemory  int `json:"consensus_node_memory"`

	NetworkEnv     string `json:"network_env"`
	BlockchainType string `json:"blockchain_type"`
	ConsensusType  string `json:"consensus_type"`
	Nodes          []Node `json:"nodes"`
	Links          []Link `json:"links"`

	SecPathMabType int     `json:"sec_path_mab_type"`
	PerLinkDelay   float64 `json:"per_link_delay"`
}

type DynamicParameters struct {
	ConsensusThreadCount int
	SecPathMabType       types.SecPathMabStrategy
	PerLinkDelay         float64
}

func NewTopology(topologyType types.TopologyType, dynamicParameters *DynamicParameters) (*Topology, error) {
	topology := &Topology{
		ConsensusThreadCount: dynamicParameters.ConsensusThreadCount,
		AccessLinkBandwidth:  8,
		ConsensusNodeCpu:     2,
		ConsensusNodeMemory:  1024,
		SecPathMabType:       int(dynamicParameters.SecPathMabType),
		PerLinkDelay:         dynamicParameters.PerLinkDelay,
	}
	err := loadInformation(topology, topologyType)
	if err != nil {
		return nil, fmt.Errorf("load information error: %v", err)
	}
	err = loadNodesAndLinks(topology, topologyType)
	if err != nil {
		return nil, fmt.Errorf("load nodes and links error: %v", err)
	}
	return topology, nil
}

func loadInformation(topology *Topology, topologyType types.TopologyType) error {
	switch topologyType {
	case types.TopologyType_HyperledgerFabric:
		topology.NetworkEnv = "fabric_test_topology"
		topology.BlockchainType = "fabric"
		topology.ConsensusType = "BFT-SMaRt"
	case types.TopologyType_FiscoBcos:
		topology.NetworkEnv = "fisco_bcos_test_topology"
		topology.BlockchainType = "fisco-bcos"
		topology.ConsensusType = "pbft"
	case types.TopologyType_ChainMaker:
		topology.NetworkEnv = "chainmaker_test_topology"
		topology.BlockchainType = "长安链"
		topology.ConsensusType = "TBFT"
	case types.TopologyType_SimplePathValidation:
		topology.NetworkEnv = "simple_path_validation"
		topology.BlockchainType = "无区块链"
		topology.ConsensusType = "无共识协议"
	case types.TopologyType_MulticastPathValidation:
		topology.NetworkEnv = "multicast_topology"
		topology.BlockchainType = "无区块链"
		topology.ConsensusType = "无共识协议"
	case types.TopologyType_SecPathMab:
		topology.NetworkEnv = "sec_path_mab_topology"
		topology.BlockchainType = "无区块链"
		topology.ConsensusType = "无共识协议"
	default:
		return fmt.Errorf("unsupported topology")
	}
	return nil
}

func loadNodesAndLinks(topology *Topology, topologyType types.TopologyType) error {
	switch topologyType {
	case types.TopologyType_HyperledgerFabric:
		result, err := file.ReadFile(TopologyPathFabric)
		if err != nil {
			return fmt.Errorf("read file error")
		}
		err = json.Unmarshal([]byte(result), &topology)
		if err != nil {
			return fmt.Errorf("unmarshal error")
		}
	case types.TopologyType_FiscoBcos:
		result, err := file.ReadFile(TopologyPathFisco)
		if err != nil {
			return fmt.Errorf("read file error")
		}
		err = json.Unmarshal([]byte(result), &topology)
		if err != nil {
			return fmt.Errorf("unmarshal error")
		}
	case types.TopologyType_ChainMaker:
		result, err := file.ReadFile(TopologyPathChainmaker)
		if err != nil {
			return fmt.Errorf("read file error")
		}
		err = json.Unmarshal([]byte(result), &topology)
		if err != nil {
			return fmt.Errorf("unmarshal error")
		}
	case types.TopologyType_SimplePathValidation:
		result, err := file.ReadFile(TopologySimplePathValidation)
		if err != nil {
			return fmt.Errorf("read file error")
		}
		err = json.Unmarshal([]byte(result), &topology)
		if err != nil {
			return fmt.Errorf("unmarshal error")
		}
	case types.TopologyType_MulticastPathValidation:
		result, err := file.ReadFile(TopologyMulticastPathValidation)
		if err != nil {
			return fmt.Errorf("read file error")
		}
		err = json.Unmarshal([]byte(result), &topology)
		if err != nil {
			return fmt.Errorf("unmarshal error")
		}
	case types.TopologyType_SecPathMab:
		result, err := file.ReadFile(TopologySecPathMab)
		if err != nil {
			return fmt.Errorf("read file error")
		}
		err = json.Unmarshal([]byte(result), &topology)
		if err != nil {
			return fmt.Errorf("unmarshal error")
		}
	}
	return nil
}
