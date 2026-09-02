package entities

import (
	"chain_simulation/configs"
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

type BlockChainParams struct {
	ConsensusThreadCount int    `json:"consensus_thread_count"`
	BlockchainType       string `json:"blockchain_type"`
	ConsensusType        string `json:"consensus_type"`
	StartDefence         bool   `json:"start_defence"`
}

type SecPathMabParams struct {
	TopologyType           int  `json:"topology_type"`
	ExperimentType         int  `json:"experiment_type"`
	SecPathMabEnabled      bool `json:"sec_path_mab_enabled"`
	SecPathMabType         int  `json:"sec_path_mab_type"`
	NumberOfHops           int  `json:"number_of_hops"`
	NumberOfSegmentsPerHop int  `json:"number_of_segments_per_hop"`
}

type TopologyParams struct {
	Nodes        []Node  `json:"nodes"`
	Links        []Link  `json:"links"`
	PerLinkDelay float64 `json:"per_link_delay"`
}

type TopologyStartParams struct {
	BlockChainParams BlockChainParams `json:"block_chain_params"`
	SecPathMabParams SecPathMabParams `json:"sec_path_mab_params"`
	TopologyParams   TopologyParams   `json:"topology_params"`
}

type DynamicParameters struct {
	ConsensusThreadCount int
	SecPathMabType       types.SecPathMabStrategy
	PerLinkDelay         float64
}

func NewTopologyStartParams(topologyType types.TopologyType, dynamicParameters *DynamicParameters) (*TopologyStartParams, error) {
	topology := &TopologyStartParams{
		BlockChainParams: BlockChainParams{
			ConsensusThreadCount: dynamicParameters.ConsensusThreadCount,
			BlockchainType:       "",
			ConsensusType:        "",
		},
		SecPathMabParams: SecPathMabParams{
			TopologyType:           configs.TopConfigInstance.SecPathMabConfig.TopologyType,
			ExperimentType:         configs.TopConfigInstance.SecPathMabConfig.ExperimentType,
			SecPathMabType:         int(dynamicParameters.SecPathMabType),
			NumberOfHops:           configs.TopConfigInstance.SecPathMabConfig.NumberOfHops,
			NumberOfSegmentsPerHop: configs.TopConfigInstance.SecPathMabConfig.NumberOfIntermediateNodes,
		},
		TopologyParams: TopologyParams{
			Nodes:        make([]Node, 0),
			Links:        make([]Link, 0),
			PerLinkDelay: dynamicParameters.PerLinkDelay,
		}}
	err := loadInformation(topology, topologyType)
	if err != nil {
		return nil, fmt.Errorf("load information error: %v", err)
	}
	err = loadNodesAndLinks(topology, topologyType)
	if err != nil {
		return nil, fmt.Errorf("load nodes and links error: %v", err)
	}
	fmt.Println("fuck")
	fmt.Println(topology)
	fmt.Println("fuck")
	return topology, nil
}

func loadInformation(topologyStartParams *TopologyStartParams, topologyType types.TopologyType) error {
	switch topologyType {
	case types.TopologyType_HyperledgerFabric:
		topologyStartParams.BlockChainParams.BlockchainType = "fabric"
		topologyStartParams.BlockChainParams.ConsensusType = "BFT-SMaRt"
	case types.TopologyType_FiscoBcos:
		topologyStartParams.BlockChainParams.BlockchainType = "fisco-bcos"
		topologyStartParams.BlockChainParams.ConsensusType = "pbft"
	case types.TopologyType_ChainMaker:
		topologyStartParams.BlockChainParams.BlockchainType = "长安链"
		topologyStartParams.BlockChainParams.ConsensusType = "TBFT"
	case types.TopologyType_SimplePathValidation:
		topologyStartParams.BlockChainParams.BlockchainType = "无区块链"
		topologyStartParams.BlockChainParams.ConsensusType = "无共识协议"
	case types.TopologyType_MulticastPathValidation:
		topologyStartParams.BlockChainParams.BlockchainType = "无区块链"
		topologyStartParams.BlockChainParams.ConsensusType = "无共识协议"
	case types.TopologyType_SecPathMab:
		topologyStartParams.BlockChainParams.BlockchainType = "无区块链"
		topologyStartParams.BlockChainParams.ConsensusType = "无共识协议"
	default:
		return fmt.Errorf("unsupported topologyStartParams")
	}
	return nil
}

func loadNodesAndLinks(topologyStartParams *TopologyStartParams, topologyType types.TopologyType) error {
	switch topologyType {
	case types.TopologyType_HyperledgerFabric:
		result, err := file.ReadFile(TopologyPathFabric)
		if err != nil {
			return fmt.Errorf("read file error")
		}
		err = json.Unmarshal([]byte(result), &topologyStartParams.TopologyParams)
		if err != nil {
			return fmt.Errorf("unmarshal error")
		}
	case types.TopologyType_FiscoBcos:
		result, err := file.ReadFile(TopologyPathFisco)
		if err != nil {
			return fmt.Errorf("read file error")
		}
		err = json.Unmarshal([]byte(result), &topologyStartParams.TopologyParams)
		if err != nil {
			return fmt.Errorf("unmarshal error")
		}
	case types.TopologyType_ChainMaker:
		result, err := file.ReadFile(TopologyPathChainmaker)
		if err != nil {
			return fmt.Errorf("read file error")
		}
		err = json.Unmarshal([]byte(result), &topologyStartParams.TopologyParams)
		if err != nil {
			return fmt.Errorf("unmarshal error")
		}
	case types.TopologyType_SimplePathValidation:
		result, err := file.ReadFile(TopologySimplePathValidation)
		if err != nil {
			return fmt.Errorf("read file error")
		}
		err = json.Unmarshal([]byte(result), &topologyStartParams.TopologyParams)
		if err != nil {
			return fmt.Errorf("unmarshal error")
		}
	case types.TopologyType_MulticastPathValidation:
		result, err := file.ReadFile(TopologyMulticastPathValidation)
		if err != nil {
			return fmt.Errorf("read file error")
		}
		err = json.Unmarshal([]byte(result), &topologyStartParams.TopologyParams)
		if err != nil {
			return fmt.Errorf("unmarshal error")
		}
	case types.TopologyType_SecPathMab:
		result, err := file.ReadFile(TopologySecPathMab)
		if err != nil {
			return fmt.Errorf("read file error")
		}
		err = json.Unmarshal([]byte(result), &topologyStartParams.TopologyParams)
		if err != nil {
			return fmt.Errorf("unmarshal error")
		}
	}
	return nil
}
