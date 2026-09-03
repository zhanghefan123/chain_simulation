package topology

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

var topologyPathMapping = map[types.TopologyType]string{
	types.TopologyType_HyperledgerFabric:       TopologyPathFabric,
	types.TopologyType_FiscoBcos:               TopologyPathFisco,
	types.TopologyType_ChainMaker:              TopologyPathChainmaker,
	types.TopologyType_SimplePathValidation:    TopologySimplePathValidation,
	types.TopologyType_MulticastPathValidation: TopologyMulticastPathValidation,
	types.TopologyType_SecPathMab:              TopologySecPathMab,
}

type BlockChainParams struct {
	ConsensusThreadCount int    `json:"consensus_thread_count"`
	BlockchainType       string `json:"blockchain_type"`
	ConsensusType        string `json:"consensus_type"`
	StartDefence         bool   `json:"start_defence"`
}

type SecPathMabParams struct {
	Enabled                bool `json:"enabled"`
	ExperimentType         int  `json:"experiment_type"`
	SecPathMabEnabled      bool `json:"sec_path_mab_enabled"`
	SecPathMabType         int  `json:"sec_path_mab_type"`
	NumberOfHops           int  `json:"number_of_hops"`
	NumberOfSegmentsPerHop int  `json:"number_of_segments_per_hop"`
}

type Params struct {
	Nodes        []Node  `json:"nodes"`
	Links        []Link  `json:"links"`
	PerLinkDelay float64 `json:"per_link_delay"`
}

type StartParams struct {
	BlockChainParams BlockChainParams `json:"block_chain_params"`
	SecPathMabParams SecPathMabParams `json:"sec_path_mab_params"`
	TopologyParams   Params           `json:"topology_params"`
}

type DynamicParameters struct {
	ConsensusThreadCount int
	SecPathMabType       types.SecPathMabStrategy
	PerLinkDelay         float64
}

func NewStartParams(topologyType types.TopologyType, dynamicParameters *DynamicParameters) (*StartParams, error) {
	topology := &StartParams{
		BlockChainParams: BlockChainParams{
			ConsensusThreadCount: dynamicParameters.ConsensusThreadCount,
			BlockchainType:       "",
			ConsensusType:        "",
		},
		SecPathMabParams: SecPathMabParams{
			Enabled:                configs.TopConfigInstance.SecPathMabConfig.Enabled,
			ExperimentType:         configs.TopConfigInstance.SecPathMabConfig.ExperimentType,
			SecPathMabType:         int(dynamicParameters.SecPathMabType),
			NumberOfHops:           configs.TopConfigInstance.SecPathMabConfig.NumberOfHops,
			NumberOfSegmentsPerHop: configs.TopConfigInstance.SecPathMabConfig.NumberOfIntermediateNodes,
		},
		TopologyParams: Params{
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

func loadInformation(topologyStartParams *StartParams, topologyType types.TopologyType) error {
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

func loadNodesAndLinks(topologyStartParams *StartParams, topologyType types.TopologyType) error {
	topologyPath, ok := topologyPathMapping[topologyType]
	if !ok {
		return fmt.Errorf("unsupported topology type: %v", topologyType)
	}

	result, err := file.ReadFile(topologyPath)
	if err != nil {
		return fmt.Errorf("read topology file %s failed: %w", topologyPath, err)
	}

	if err := json.Unmarshal([]byte(result), &topologyStartParams.TopologyParams); err != nil {
		return fmt.Errorf("unmarshal topology file %s failed: %w", topologyPath, err)
	}

	return nil
}
