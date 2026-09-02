package sec_path_mab_topology_generator

import (
	"chain_simulation/entities"
	"chain_simulation/entities/sec_path_mab"
	"chain_simulation/utils/file"
	"encoding/json"
	"fmt"
)

type NonLinearTopologyDescription struct {
	NumberOfHops              int                            `json:"number_of_hops"`
	NumberOfIntermediateNodes int                            `json:"number_of_intermediate_nodes"`
	SourceDestParams          sec_path_mab.SourceDestParams  `json:"source_dest_params"`
	Nodes                     []sec_path_mab.SimNodeParam    `json:"nodes"`
	PvLinks                   []sec_path_mab.SimAbsLinkParam `json:"pv_links"`
	Links                     []sec_path_mab.SimLinkParam    `json:"links"`
}

// GenerateNonLinearTopologyDescription 进行拓扑的生成
func GenerateNonLinearTopologyDescription(numberOfHops, numberOfIntermediateNodes int, lowRatio, highRatio float64) *NonLinearTopologyDescription {
	finalResult := &NonLinearTopologyDescription{
		NumberOfHops:              numberOfHops,
		NumberOfIntermediateNodes: numberOfIntermediateNodes,
	}
	finalResult.FillSourceDestParams()
	indexToNodeParamMapping := finalResult.FillNodes(lowRatio, highRatio)
	finalResult.FillPvLinks(indexToNodeParamMapping)
	finalResult.FillRealLinks(indexToNodeParamMapping)
	return finalResult
}

// MarshalNonLinearTopologyDescription 将拓扑描述放到文件之中
func (osmdTopologyDescription *NonLinearTopologyDescription) MarshalNonLinearTopologyDescription(destinationFile string) {
	bytes, err := json.MarshalIndent(osmdTopologyDescription, "", " ")
	if err != nil {
		return
	}
	err = file.WriteStringIntoFile(destinationFile, string(bytes))
	if err != nil {
		return
	}
}

// FillSourceDestParams 填充 topologyDescription 之中的 source dest params
func (osmdTopologyDescription *NonLinearTopologyDescription) FillSourceDestParams() {
	osmdTopologyDescription.SourceDestParams = sec_path_mab.SourceDestParams{
		Source:      "EndHost-1",
		Destination: fmt.Sprintf("EndHost-%d", 1+osmdTopologyDescription.NumberOfHops*(osmdTopologyDescription.NumberOfIntermediateNodes+1)),
	}
}

// FillNodes 填充节点
func (osmdTopologyDescription *NonLinearTopologyDescription) FillNodes(lowRatio, highRatio float64) map[int]sec_path_mab.SimNodeParam {
	indexToNodeMapping := map[int]sec_path_mab.SimNodeParam{}
	if osmdTopologyDescription.NumberOfIntermediateNodes <= 0 {
		osmdTopologyDescription.NumberOfIntermediateNodes = 2
	}
	numberOfNodes := 1 + osmdTopologyDescription.NumberOfHops*(osmdTopologyDescription.NumberOfIntermediateNodes+1)
	osmdTopologyDescription.Nodes = make([]sec_path_mab.SimNodeParam, numberOfNodes)
	osmdTopologyDescription.Nodes[0] = sec_path_mab.SimNodeParam{
		Index: 1,
		Type:  "EndHost",
	}
	indexToNodeMapping[1] = osmdTopologyDescription.Nodes[0]
	currentHop := 0
	index := 1
	for {
		for i := 0; i < osmdTopologyDescription.NumberOfIntermediateNodes; i++ {
			ratio := lowRatio
			if i%2 == 1 {
				ratio = highRatio
			}
			osmdTopologyDescription.Nodes[index] = sec_path_mab.SimNodeParam{
				Index: index + 1,
				Type:  "NormalRouter",
				CorruptRatio: entities.RatioDistribution{
					Start: ratio,
					End:   ratio,
				},
				CorruptSpecialPacketRatio: entities.RatioDistribution{
					Start: 0,
					End:   0,
				},
			}
			indexToNodeMapping[index+1] = osmdTopologyDescription.Nodes[index]
			index += 1
		}

		if index == (numberOfNodes - 1) {
			osmdTopologyDescription.Nodes[index] = sec_path_mab.SimNodeParam{
				Index: index + 1,
				Type:  "EndHost",
			}
		} else {
			osmdTopologyDescription.Nodes[index] = sec_path_mab.SimNodeParam{
				Index: index + 1,
				Type:  "PathValidationRouter",
			}
		}
		indexToNodeMapping[index+1] = osmdTopologyDescription.Nodes[index]
		index += 1
		currentHop += 1

		if currentHop >= osmdTopologyDescription.NumberOfHops {
			break
		}
	}
	return indexToNodeMapping
}

// FillPvLinks 进行抽象边的填充
func (osmdTopologyDescription *NonLinearTopologyDescription) FillPvLinks(indexToNodeParamMapping map[int]sec_path_mab.SimNodeParam) {
	numberOfPvLinks := osmdTopologyDescription.NumberOfHops * osmdTopologyDescription.NumberOfIntermediateNodes
	osmdTopologyDescription.PvLinks = make([]sec_path_mab.SimAbsLinkParam, numberOfPvLinks)

	currentHop := 0
	currentNodeIndex := 0
	pvLinkIndex := 0
	for {
		source := indexToNodeParamMapping[currentNodeIndex+1]
		target := indexToNodeParamMapping[currentNodeIndex+osmdTopologyDescription.NumberOfIntermediateNodes+2]
		for i := 0; i < osmdTopologyDescription.NumberOfIntermediateNodes; i++ {
			osmdTopologyDescription.PvLinks[pvLinkIndex] = sec_path_mab.SimAbsLinkParam{
				SourceNode:       source,
				TargetNode:       target,
				IntermediateNode: indexToNodeParamMapping[currentNodeIndex+2+i],
			}
			pvLinkIndex += 1
		}
		currentNodeIndex += osmdTopologyDescription.NumberOfIntermediateNodes + 1
		currentHop += 1

		if currentHop >= osmdTopologyDescription.NumberOfHops {
			break
		}
	}
}

func (osmdTopologyDescription *NonLinearTopologyDescription) FillRealLinks(indexToNodeParamMapping map[int]sec_path_mab.SimNodeParam) {
	numberOfRealLinks := osmdTopologyDescription.NumberOfHops * osmdTopologyDescription.NumberOfIntermediateNodes * 2
	osmdTopologyDescription.Links = make([]sec_path_mab.SimLinkParam, numberOfRealLinks)

	currentHop := 0
	currentNodeIndex := 0
	realLinkIndex := 0
	for {
		source := indexToNodeParamMapping[currentNodeIndex+1]
		target := indexToNodeParamMapping[currentNodeIndex+osmdTopologyDescription.NumberOfIntermediateNodes+2]
		for i := 0; i < osmdTopologyDescription.NumberOfIntermediateNodes; i++ {
			intermediate := indexToNodeParamMapping[currentNodeIndex+2+i]
			osmdTopologyDescription.Links[realLinkIndex] = sec_path_mab.SimLinkParam{
				SourceNode: source,
				TargetNode: intermediate,
			}
			realLinkIndex += 1

			osmdTopologyDescription.Links[realLinkIndex] = sec_path_mab.SimLinkParam{
				SourceNode: intermediate,
				TargetNode: target,
			}
			realLinkIndex += 1
		}
		currentNodeIndex += osmdTopologyDescription.NumberOfIntermediateNodes + 1
		currentHop += 1

		if currentHop >= osmdTopologyDescription.NumberOfHops {
			break
		}

	}
}
