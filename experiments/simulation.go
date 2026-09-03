package experiments

import (
	"chain_simulation/modules/experiment_related/scheduler"
	"errors"
	"fmt"

	"chain_simulation/configs"
	"chain_simulation/entities"
	"chain_simulation/modules/backend_manager"
	"chain_simulation/utils/file"
)

var experimentIndex = 0

// SingleSimulation runs one complete experiment lifecycle.
func SingleSimulation(configurationSetting *entities.ConfigurationSetting, events []*entities.Event) error {
	if configurationSetting == nil {
		return fmt.Errorf("configuration setting must not be nil")
	}

	if err := file.ModifyYml(
		configs.TopConfigInstance.PathConfig.ConfigurationYml,
		configurationSetting.Mapping,
	); err != nil {
		return fmt.Errorf("modify simulation configuration: %w", err)
	}

	if err := backend_manager.StartBackendService(experimentIndex); err != nil {
		return fmt.Errorf("start backend for experiment %d: %w", experimentIndex, err)
	}

	runErr := scheduler.Run(events)
	stopErr := backend_manager.StopBackendService()
	if err := errors.Join(runErr, stopErr); err != nil {
		return fmt.Errorf("run simulation for mapping %v: %w", configurationSetting.Mapping, err)
	}

	fmt.Printf("simulation for mapping %v is finished\n", configurationSetting.Mapping)
	experimentIndex++
	return nil
}
