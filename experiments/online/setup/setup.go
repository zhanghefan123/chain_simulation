package setup

import (
	"chain_simulation/modules/experiment_related/sec_path_mab_topology_generator"
	"fmt"
	"path"

	"chain_simulation/configs"
	onlineexecutor "chain_simulation/experiments/online/executor"
)

const simulatorTopologyDirectory = "/home/zhf/Projects/emulator/backend/resources/online_topologies"

// InitOnlineExperimentEnvironment loads configuration, refreshes runtime
// candidates, and writes the topology descriptions required by the simulator
// and backend.
func InitOnlineExperimentEnvironment() error {
	if err := configs.InitTopConfig(); err != nil {
		return fmt.Errorf("initialize top-level config: %w", err)
	}
	onlineexecutor.RefreshMaliciousCandidatesFromConfig()
	return UpdateExperimentEnvironment()
}

// UpdateExperimentEnvironment regenerates topology files from the current
// SecPathMAB configuration.
func UpdateExperimentEnvironment() error {
	simulatorTopologyPath := path.Join(simulatorTopologyDirectory, "sec_path_mab_topology.json")
	backendTopologyPath := path.Join(
		configs.TopConfigInstance.PathConfig.ResourcesPath,
		"topologies/sec_path_mab_topology.json",
	)

	if configs.TopConfigInstance.SecPathMabConfig.Enabled {
		description := sec_path_mab_topology_generator.GenerateNonLinearTopologyDescription(
			configs.TopConfigInstance.SecPathMabConfig.NumberOfHops,
			configs.TopConfigInstance.SecPathMabConfig.NumberOfIntermediateNodes,
			configs.TopConfigInstance.SecPathMabConfig.LowCorruptRatio,
			configs.TopConfigInstance.SecPathMabConfig.HighCorruptRatio,
		)
		if err := description.WriteNonLinearTopologyDescription(simulatorTopologyPath); err != nil {
			return err
		}
		backendDescription := sec_path_mab_topology_generator.GenerateBuildTopologyDescription(description)
		return backendDescription.WriteBuildTopologyDescription(backendTopologyPath)
	} else {
		description := sec_path_mab_topology_generator.GenerateLinearTopologyDescription(
			configs.TopConfigInstance.SecPathMabConfig.NumberOfHops,
		)
		if err := description.WriteLinearTopologyDescription(simulatorTopologyPath); err != nil {
			return err
		}
		backendDescription := sec_path_mab_topology_generator.GenerateBuildTopologyDescription(
			description.ToOsmdTopologyDescription(),
		)
		return backendDescription.WriteBuildTopologyDescription(backendTopologyPath)
	}
}
