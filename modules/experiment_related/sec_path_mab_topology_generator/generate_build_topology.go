package sec_path_mab_topology_generator

import (
	"encoding/json"
	"fmt"

	"chain_simulation/entities"
	"chain_simulation/utils/file"
)

type BuildTopologyDescription struct {
	Nodes []entities.Node `json:"nodes"`
	Links []entities.Link `json:"links"`
}

func GenerateBuildTopologyDescription(osmdTopologyDescription *NonLinearTopologyDescription) *BuildTopologyDescription {
	finalResult := &BuildTopologyDescription{}
	indexToBuildNodeMapping := finalResult.FillNodes(osmdTopologyDescription)
	finalResult.FillLinks(osmdTopologyDescription, indexToBuildNodeMapping)
	return finalResult
}

// WriteBuildTopologyDescription writes the backend topology description as JSON.
func (buildTopologyDescription *BuildTopologyDescription) WriteBuildTopologyDescription(destinationFile string) error {
	bytes, err := json.MarshalIndent(buildTopologyDescription, "", " ")
	if err != nil {
		return fmt.Errorf("marshal build topology description: %w", err)
	}
	if err := file.WriteStringIntoFile(destinationFile, string(bytes)); err != nil {
		return fmt.Errorf("write build topology description to %s: %w", destinationFile, err)
	}
	return nil
}

func (buildTopologyDescription *BuildTopologyDescription) FillNodes(osmdTopologyDescription *NonLinearTopologyDescription) map[int]entities.Node {
	indexToBuildNodeMapping := make(map[int]entities.Node)
	buildTopologyDescription.Nodes = make([]entities.Node, len(osmdTopologyDescription.Nodes))
	index := 0
	for osmdNodeIndex := range osmdTopologyDescription.Nodes {
		innerRouterType := -1
		osmdNode := osmdTopologyDescription.Nodes[osmdNodeIndex]
		if osmdNode.Type == "EndHost" || osmdNode.Type == "PathValidationRouter" {
			innerRouterType = 2
		} else {
			innerRouterType = 1
		}

		buildNode := entities.Node{
			Index: osmdNode.Index,
			Type:  "LirNode",
			X:     10.0 * float64(osmdNode.Index),
			Y:     10.0 * float64(osmdNode.Index),
			SpecialParmas: entities.SpecialParams{
				InnerRouterType: innerRouterType,
				MaliciousParams: entities.MaliciousParams{
					CorruptRatio:        osmdNode.CorruptRatio,
					CorruptSpecialRatio: osmdNode.CorruptSpecialPacketRatio,
				},
			},
		}
		indexToBuildNodeMapping[osmdNode.Index] = buildNode
		buildTopologyDescription.Nodes[index] = buildNode
		index += 1
	}
	return indexToBuildNodeMapping
}

func (buildTopologyDescription *BuildTopologyDescription) FillLinks(osmdTopologyDescription *NonLinearTopologyDescription, indexToBuildNodeMapping map[int]entities.Node) {
	buildTopologyDescription.Links = make([]entities.Link, len(osmdTopologyDescription.Links))
	for realLinkIndex := range osmdTopologyDescription.Links {
		realLink := osmdTopologyDescription.Links[realLinkIndex]
		buildTopologyDescription.Links[realLinkIndex] = entities.Link{
			SourceNode: indexToBuildNodeMapping[realLink.SourceNode.Index],
			TargetNode: indexToBuildNodeMapping[realLink.TargetNode.Index],
			LinkType:   "backbone",
		}
	}
}
