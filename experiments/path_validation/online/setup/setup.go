package setup

import (
	"chain_simulation/configs"
	"chain_simulation/entities/types"
	onlineexecutor "chain_simulation/experiments/path_validation/online/executor"
	"chain_simulation/modules/sec_path_mab_topology_generator"
	"fmt"
	"path"
)

func UpdateExperimentEnvironment() error {
	if configs.TopConfigInstance.SecPathMabConfig.TopologyType == (int)(types.SecPathMabTopologyType_NON_LINEAR_TEST_TOPOLOGY) {
		// 生成结构
		osmdTopologyDesc := sec_path_mab_topology_generator.GenerateNonLinearTopologyDescription(
			configs.TopConfigInstance.SecPathMabConfig.NumberOfHops,
			configs.TopConfigInstance.SecPathMabConfig.NumberOfIntermediateNodes,
			configs.TopConfigInstance.SecPathMabConfig.LowCorruptRatio,
			configs.TopConfigInstance.SecPathMabConfig.HighCorruptRatio)
		destinationFilePath := fmt.Sprintf("%s/%s", "/home/zhf/Projects/emulator/backend/resources/online_topologies",
			"sec_path_mab_topology.json")
		osmdTopologyDesc.MarshalNonLinearTopologyDescription(destinationFilePath)
		// 生成描述
		buildTopologyDescription := sec_path_mab_topology_generator.GenerateBuildTopologyDescription(osmdTopologyDesc)
		// 输出
		outputFilePath := path.Join(configs.TopConfigInstance.PathConfig.ResourcesPath, "/topologies/sec_path_mab_topology.json")
		buildTopologyDescription.MarshalBuildTopologyDescription(outputFilePath)
	} else if configs.TopConfigInstance.SecPathMabConfig.TopologyType == (int)(types.SecPathMabTopologyType_LINEAR_TEST_TOPOLOGY) {
		linearTopologyDesc := sec_path_mab_topology_generator.GenerateLinearTopologyDescription(
			configs.TopConfigInstance.SecPathMabConfig.NumberOfHops)
		destinationFilePath := fmt.Sprintf("%s/%s", "/home/zhf/Projects/emulator/backend/resources/online_topologies",
			"sec_path_mab_topology.json")
		linearTopologyDesc.MarshalLinearTopologyDescription(destinationFilePath)
		// 生成描述
		buildTopologyDescription := sec_path_mab_topology_generator.GenerateBuildTopologyDescription(linearTopologyDesc.ToOsmdTopologyDescription())
		// 输出
		outputFilePath := path.Join(configs.TopConfigInstance.PathConfig.ResourcesPath, "/topologies/sec_path_mab_topology.json")
		buildTopologyDescription.MarshalBuildTopologyDescription(outputFilePath)
	}
	return nil
}

func InitOnlineExperimentEnvironment() error {
	if err := configs.InitTopConfig(); err != nil {
		return fmt.Errorf("init top config failed: %w", err)
	}
	onlineexecutor.RefreshMaliciousCandidatesFromConfig()

	//osmdTopologyDesc := sec_path_mab_topology_generator.GenerateNonLinearTopologyDescription(
	//	configs.TopConfigInstance.SecPathMabConfig.NumberOfHops,
	//	configs.TopConfigInstance.SecPathMabConfig.NumberOfIntermediateNodes,
	//	configs.TopConfigInstance.SecPathMabConfig.LowCorruptRatio,
	//	configs.TopConfigInstance.SecPathMabConfig.HighCorruptRatio,
	//)
	//destinationFilePath := fmt.Sprintf("%s/%s",
	//	"/home/zhf/Projects/emulator/backend/resources/online_topologies",
	//	"sec_path_mab_topology.json")
	//osmdTopologyDesc.MarshalNonLinearTopologyDescription(destinationFilePath)
	//
	//buildTopologyDescription := sec_path_mab_topology_generator.GenerateBuildTopologyDescription(osmdTopologyDesc)
	//outputFilePath := path.Join(configs.TopConfigInstance.PathConfig.ResourcesPath, "/topologies/sec_path_mab_topology.json")
	//buildTopologyDescription.MarshalBuildTopologyDescription(outputFilePath)

	if configs.TopConfigInstance.SecPathMabConfig.TopologyType == (int)(types.SecPathMabTopologyType_NON_LINEAR_TEST_TOPOLOGY) {
		// 生成结构
		osmdTopologyDesc := sec_path_mab_topology_generator.GenerateNonLinearTopologyDescription(
			configs.TopConfigInstance.SecPathMabConfig.NumberOfHops,
			configs.TopConfigInstance.SecPathMabConfig.NumberOfIntermediateNodes,
			configs.TopConfigInstance.SecPathMabConfig.LowCorruptRatio,
			configs.TopConfigInstance.SecPathMabConfig.HighCorruptRatio)
		destinationFilePath := fmt.Sprintf("%s/%s", "/home/zhf/Projects/emulator/backend/resources/online_topologies",
			"sec_path_mab_topology.json")
		osmdTopologyDesc.MarshalNonLinearTopologyDescription(destinationFilePath)
		// 生成描述
		buildTopologyDescription := sec_path_mab_topology_generator.GenerateBuildTopologyDescription(osmdTopologyDesc)
		// 输出
		outputFilePath := path.Join(configs.TopConfigInstance.PathConfig.ResourcesPath, "/topologies/sec_path_mab_topology.json")
		buildTopologyDescription.MarshalBuildTopologyDescription(outputFilePath)
	} else if configs.TopConfigInstance.SecPathMabConfig.TopologyType == (int)(types.SecPathMabTopologyType_LINEAR_TEST_TOPOLOGY) {
		linearTopologyDesc := sec_path_mab_topology_generator.GenerateLinearTopologyDescription(
			configs.TopConfigInstance.SecPathMabConfig.NumberOfHops)
		destinationFilePath := fmt.Sprintf("%s/%s", "/home/zhf/Projects/emulator/backend/resources/online_topologies",
			"sec_path_mab_topology.json")
		linearTopologyDesc.MarshalLinearTopologyDescription(destinationFilePath)
		// 生成描述
		buildTopologyDescription := sec_path_mab_topology_generator.GenerateBuildTopologyDescription(linearTopologyDesc.ToOsmdTopologyDescription())
		// 输出
		outputFilePath := path.Join(configs.TopConfigInstance.PathConfig.ResourcesPath, "/topologies/sec_path_mab_topology.json")
		buildTopologyDescription.MarshalBuildTopologyDescription(outputFilePath)
	}
	return nil
}
