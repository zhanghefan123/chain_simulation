package entities

import (
	"chain_simulation/configs"
	"chain_simulation/entities/types"
	"chain_simulation/utils/file"
	"encoding/json"
	"fmt"
)

var (
	TopologyPathFabric     = "../resources/topologies/fabric_test_topology.json"
	TopologyPathFisco      = "../resources/topologies/fisco_bcos_test_topology.json"
	TopologyPathChainmaker = "../resources/topologies/chainmaker_test_topology.json"
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
}

func NewTopology(topologyType types.TopologyType) (*Topology, error) {
	topology := &Topology{
		ConsensusThreadCount: configs.TopConfigInstance.ConsensusConfig.ThreadCount,
		AccessLinkBandwidth:  8,
		ConsensusNodeCpu:     2,
		ConsensusNodeMemory:  1024,
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
		topology.NetworkEnv = types.TopologyType_HyperledgerFabric.String()
		topology.BlockchainType = "fabric"
		topology.ConsensusType = "BFT-SMaRt"
	case types.TopologyType_FiscoBcos:
		topology.NetworkEnv = types.TopologyType_FiscoBcos.String()
		topology.BlockchainType = "fisco-bcos"
		topology.ConsensusType = "pbft"
	case types.TopologyType_ChainMaker:
		topology.NetworkEnv = types.TopologyType_ChainMaker.String()
		topology.BlockchainType = "长安链"
		topology.ConsensusType = "TBFT"
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
	}
	return nil
}
