package setup

import (
	"chain_simulation/configs"
	"chain_simulation/modules/sec_path_mab_topology_generator"
	"fmt"
	"path"
)

func InitOnlineExperimentEnvironment() error {
	if err := configs.InitTopConfig(); err != nil {
		return fmt.Errorf("init top config failed: %w", err)
	}

	osmdTopologyDesc := sec_path_mab_topology_generator.GenerateOsmdTopologyDescription(
		configs.TopConfigInstance.SecPathMabConfig.NumberOfHops,
		configs.TopConfigInstance.SecPathMabConfig.LowCorruptRatio,
		configs.TopConfigInstance.SecPathMabConfig.HighCorruptRatio,
	)
	destinationFilePath := fmt.Sprintf("%s/%s",
		"/home/zhf/Projects/emulator/backend/resources/online_topologies",
		"sec_path_mab_topology.json")
	osmdTopologyDesc.MarshalOsmdTopologyDescription(destinationFilePath)

	buildTopologyDescription := sec_path_mab_topology_generator.GenerateBuildTopologyDescription(osmdTopologyDesc)
	outputFilePath := path.Join(configs.TopConfigInstance.PathConfig.ResourcesPath, "/topologies/sec_path_mab_topology.json")
	buildTopologyDescription.MarshalBuildTopologyDescription(outputFilePath)
	return nil
}
