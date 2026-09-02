package sec_path_mab_topology_generator

import (
	"chain_simulation/entities/sec_path_mab"
	"chain_simulation/utils/file"
	"encoding/json"
	"fmt"
)

// LinearTopologyDescription 线性拓扑描述，用于 throughput 测试。
// 仅考虑 number_of_hops，PathValidationRouter 之间不插入 NormalRouter。
type LinearTopologyDescription struct {
	NumberOfHops     int                            `json:"number_of_hops"`
	SourceDestParams sec_path_mab.SourceDestParams  `json:"source_dest_params"`
	Nodes            []sec_path_mab.SimNodeParam    `json:"nodes"`
	PvLinks          []sec_path_mab.SimAbsLinkParam `json:"pv_links"`
	Links            []sec_path_mab.SimLinkParam    `json:"links"`
}

// GenerateLinearTopologyDescription 生成线性拓扑：EndHost -> PV -> ... -> EndHost
func GenerateLinearTopologyDescription(numberOfHops int) *LinearTopologyDescription {
	finalResult := &LinearTopologyDescription{
		NumberOfHops: numberOfHops,
	}
	finalResult.FillSourceDestParams()
	indexToNodeParamMapping := finalResult.FillNodes()
	finalResult.FillPvLinks(indexToNodeParamMapping)
	finalResult.FillRealLinks(indexToNodeParamMapping)
	return finalResult
}

// MarshalLinearTopologyDescription 将拓扑描述写入文件
func (linearTopologyDescription *LinearTopologyDescription) MarshalLinearTopologyDescription(destinationFile string) {
	bytes, err := json.MarshalIndent(linearTopologyDescription, "", " ")
	if err != nil {
		return
	}
	err = file.WriteStringIntoFile(destinationFile, string(bytes))
	if err != nil {
		return
	}
}

// ToOsmdTopologyDescription 转换为 NonLinearTopologyDescription，便于复用 build topology 生成逻辑
func (linearTopologyDescription *LinearTopologyDescription) ToOsmdTopologyDescription() *NonLinearTopologyDescription {
	return &NonLinearTopologyDescription{
		NumberOfHops:     linearTopologyDescription.NumberOfHops,
		SourceDestParams: linearTopologyDescription.SourceDestParams,
		Nodes:            linearTopologyDescription.Nodes,
		PvLinks:          linearTopologyDescription.PvLinks,
		Links:            linearTopologyDescription.Links,
	}
}

func (linearTopologyDescription *LinearTopologyDescription) FillSourceDestParams() {
	linearTopologyDescription.SourceDestParams = sec_path_mab.SourceDestParams{
		Source:      "EndHost-1",
		Destination: fmt.Sprintf("EndHost-%d", 1+linearTopologyDescription.NumberOfHops),
	}
}

func (linearTopologyDescription *LinearTopologyDescription) FillNodes() map[int]sec_path_mab.SimNodeParam {
	indexToNodeMapping := map[int]sec_path_mab.SimNodeParam{}
	numberOfNodes := 1 + linearTopologyDescription.NumberOfHops
	linearTopologyDescription.Nodes = make([]sec_path_mab.SimNodeParam, numberOfNodes)
	linearTopologyDescription.Nodes[0] = sec_path_mab.SimNodeParam{
		Index: 1,
		Type:  "EndHost",
	}
	indexToNodeMapping[1] = linearTopologyDescription.Nodes[0]

	currentHop := 0
	index := 1
	for {
		if index == numberOfNodes-1 {
			linearTopologyDescription.Nodes[index] = sec_path_mab.SimNodeParam{
				Index: index + 1,
				Type:  "EndHost",
			}
		} else {
			linearTopologyDescription.Nodes[index] = sec_path_mab.SimNodeParam{
				Index: index + 1,
				Type:  "PathValidationRouter",
			}
		}
		indexToNodeMapping[index+1] = linearTopologyDescription.Nodes[index]
		index++
		currentHop++

		if currentHop >= linearTopologyDescription.NumberOfHops {
			break
		}
	}
	return indexToNodeMapping
}

func (linearTopologyDescription *LinearTopologyDescription) FillPvLinks(indexToNodeParamMapping map[int]sec_path_mab.SimNodeParam) {
	linearTopologyDescription.PvLinks = make([]sec_path_mab.SimAbsLinkParam, linearTopologyDescription.NumberOfHops)
	for i := 0; i < linearTopologyDescription.NumberOfHops; i++ {
		linearTopologyDescription.PvLinks[i] = sec_path_mab.SimAbsLinkParam{
			SourceNode:       indexToNodeParamMapping[i+1],
			TargetNode:       indexToNodeParamMapping[i+2],
			IntermediateNode: sec_path_mab.SimNodeParam{},
		}
	}
}

func (linearTopologyDescription *LinearTopologyDescription) FillRealLinks(indexToNodeParamMapping map[int]sec_path_mab.SimNodeParam) {
	linearTopologyDescription.Links = make([]sec_path_mab.SimLinkParam, linearTopologyDescription.NumberOfHops)
	for i := 0; i < linearTopologyDescription.NumberOfHops; i++ {
		linearTopologyDescription.Links[i] = sec_path_mab.SimLinkParam{
			SourceNode: indexToNodeParamMapping[i+1],
			TargetNode: indexToNodeParamMapping[i+2],
		}
	}
}
