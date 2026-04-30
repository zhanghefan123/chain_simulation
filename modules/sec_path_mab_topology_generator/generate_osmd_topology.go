package sec_path_mab_topology_generator

import (
	"chain_simulation/entities"
	"chain_simulation/entities/sec_path_mab"
	"chain_simulation/utils/file"
	"encoding/json"
	"fmt"
)

type OsmdTopologyDescription struct {
	NumberOfHops     int                            `json:"number_of_hops"`
	SourceDestParams sec_path_mab.SourceDestParams  `json:"source_dest_params"`
	Nodes            []sec_path_mab.SimNodeParam    `json:"nodes"`
	PvLinks          []sec_path_mab.SimAbsLinkParam `json:"pv_links"`
	Links            []sec_path_mab.SimLinkParam    `json:"links"`
}

// GenerateOsmdTopologyDescription 进行拓扑的生成
func GenerateOsmdTopologyDescription(numberOfHops int, lowRatio, highRatio float64) *OsmdTopologyDescription {
	finalResult := &OsmdTopologyDescription{
		NumberOfHops: numberOfHops,
	}
	finalResult.FillSourceDestParams()
	indexToNodeParamMapping := finalResult.FillNodes(lowRatio, highRatio)
	finalResult.FillPvLinks(indexToNodeParamMapping)
	finalResult.FillRealLinks(indexToNodeParamMapping)
	return finalResult
}

// MarshalOsmdTopologyDescription 将拓扑描述放到文件之中
func (osmdTopologyDescription *OsmdTopologyDescription) MarshalOsmdTopologyDescription(destinationFile string) {
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
func (osmdTopologyDescription *OsmdTopologyDescription) FillSourceDestParams() {
	osmdTopologyDescription.SourceDestParams = sec_path_mab.SourceDestParams{
		Source:      "EndHost-1",
		Destination: fmt.Sprintf("EndHost-%d", 1+osmdTopologyDescription.NumberOfHops*3),
	}
}

// FillNodes 填充节点
func (osmdTopologyDescription *OsmdTopologyDescription) FillNodes(lowRatio, highRatio float64) map[int]sec_path_mab.SimNodeParam {
	indexToNodeMapping := map[int]sec_path_mab.SimNodeParam{}
	numberOfNodes := 1 + osmdTopologyDescription.NumberOfHops*3
	osmdTopologyDescription.Nodes = make([]sec_path_mab.SimNodeParam, numberOfNodes)
	osmdTopologyDescription.Nodes[0] = sec_path_mab.SimNodeParam{
		Index: 1,
		Type:  "EndHost",
	}
	indexToNodeMapping[1] = osmdTopologyDescription.Nodes[0]
	currentHop := 0
	index := 1
	for {
		osmdTopologyDescription.Nodes[index] = sec_path_mab.SimNodeParam{
			Index: index + 1,
			Type:  "NormalRouter",
			CorruptRatio: entities.RatioDistribution{
				Start: lowRatio,
				End:   lowRatio,
			},
			CorruptSpecialPacketRatio: entities.RatioDistribution{
				Start: 0,
				End:   0,
			},
		}
		indexToNodeMapping[index+1] = osmdTopologyDescription.Nodes[index]

		index += 1
		osmdTopologyDescription.Nodes[index] = sec_path_mab.SimNodeParam{
			Index: index + 1,
			Type:  "NormalRouter",
			CorruptRatio: entities.RatioDistribution{
				Start: highRatio,
				End:   highRatio,
			},
			CorruptSpecialPacketRatio: entities.RatioDistribution{
				Start: 0,
				End:   0,
			},
		}
		indexToNodeMapping[index+1] = osmdTopologyDescription.Nodes[index]

		index += 1
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
func (osmdTopologyDescription *OsmdTopologyDescription) FillPvLinks(indexToNodeParamMapping map[int]sec_path_mab.SimNodeParam) {
	numberOfPvLinks := osmdTopologyDescription.NumberOfHops * 2
	osmdTopologyDescription.PvLinks = make([]sec_path_mab.SimAbsLinkParam, numberOfPvLinks)

	currentHop := 0
	currentNodeIndex := 0
	pvLinkIndex := 0
	for {
		osmdTopologyDescription.PvLinks[pvLinkIndex] = sec_path_mab.SimAbsLinkParam{
			SourceNode:       indexToNodeParamMapping[currentNodeIndex+1],
			TargetNode:       indexToNodeParamMapping[currentNodeIndex+4],
			IntermediateNode: indexToNodeParamMapping[currentNodeIndex+2],
		}
		pvLinkIndex += 1

		osmdTopologyDescription.PvLinks[pvLinkIndex] = sec_path_mab.SimAbsLinkParam{
			SourceNode:       indexToNodeParamMapping[currentNodeIndex+1],
			TargetNode:       indexToNodeParamMapping[currentNodeIndex+4],
			IntermediateNode: indexToNodeParamMapping[currentNodeIndex+3],
		}
		pvLinkIndex += 1
		currentNodeIndex += 3
		currentHop += 1

		if currentHop >= osmdTopologyDescription.NumberOfHops {
			break
		}
	}
}

func (osmdTopologyDescription *OsmdTopologyDescription) FillRealLinks(indexToNodeParamMapping map[int]sec_path_mab.SimNodeParam) {
	numberOfRealLinks := osmdTopologyDescription.NumberOfHops * 4
	osmdTopologyDescription.Links = make([]sec_path_mab.SimLinkParam, numberOfRealLinks)

	currentHop := 0
	currentNodeIndex := 0
	realLinkIndex := 0
	for {
		osmdTopologyDescription.Links[realLinkIndex] = sec_path_mab.SimLinkParam{
			SourceNode: indexToNodeParamMapping[currentNodeIndex+1],
			TargetNode: indexToNodeParamMapping[currentNodeIndex+2],
		}
		realLinkIndex += 1

		osmdTopologyDescription.Links[realLinkIndex] = sec_path_mab.SimLinkParam{
			SourceNode: indexToNodeParamMapping[currentNodeIndex+2],
			TargetNode: indexToNodeParamMapping[currentNodeIndex+4],
		}
		realLinkIndex += 1

		osmdTopologyDescription.Links[realLinkIndex] = sec_path_mab.SimLinkParam{
			SourceNode: indexToNodeParamMapping[currentNodeIndex+1],
			TargetNode: indexToNodeParamMapping[currentNodeIndex+3],
		}
		realLinkIndex += 1

		osmdTopologyDescription.Links[realLinkIndex] = sec_path_mab.SimLinkParam{
			SourceNode: indexToNodeParamMapping[currentNodeIndex+3],
			TargetNode: indexToNodeParamMapping[currentNodeIndex+4],
		}
		realLinkIndex += 1
		currentNodeIndex += 3
		currentHop += 1

		if currentHop >= osmdTopologyDescription.NumberOfHops {
			break
		}

	}
}
