package executor

import (
	"chain_simulation/entities"
	"chain_simulation/experiments"
	onlineconfig "chain_simulation/experiments/online/config"
	"fmt"
)

type DifferentBatchSizeEventGenerator func(currentExperimentIndex int, setting *entities.ConfigurationSetting) ([]*entities.Event, error)

func RunDifferentBatchSizeExperiments(
	scenarioLabel string,
	configurationSettings []*entities.ConfigurationSetting,
	generateEvents DifferentBatchSizeEventGenerator,
) error {
	for index, configurationSetting := range configurationSettings {
		experimentIndex := index + 1
		for runIndex := 1; runIndex <= onlineconfig.ExperimentRunCount(); runIndex++ {
			settingForRun := onlineconfig.ConfigurationSettingForRun(configurationSetting, runIndex)
			resultDir := onlineconfig.ResultDirForSetting(settingForRun, experimentIndex)
			if onlineconfig.IsResultExists(resultDir) {
				fmt.Printf("skip existing result for %s (run %d/%d): %s\n",
					scenarioLabel, runIndex, onlineconfig.ExperimentRunCount(), resultDir)
				continue
			}

			fmt.Printf("run sec path mab experiment for %s (run %d/%d): %s\n",
				scenarioLabel, runIndex, onlineconfig.ExperimentRunCount(), resultDir)
			secPathMabEvents, err := generateEvents(experimentIndex, settingForRun)
			fmt.Printf("number of events: %d\n", len(secPathMabEvents))
			if err != nil {
				return fmt.Errorf("generate events for %s run %d failed: %w", scenarioLabel, runIndex, err)
			}

			err = experiments.SingleSimulation(settingForRun, secPathMabEvents)
			if err != nil {
				return fmt.Errorf("sec path mab batch experiment failed for %s run %d: %w", scenarioLabel, runIndex, err)
			}
		}
	}
	return nil
}
