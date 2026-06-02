package config

import (
	"chain_simulation/entities"
	"chain_simulation/utils/dir"
	"fmt"
	"os"
	"strconv"
)

const ResultBaseDir = "/home/zhf/Projects/emulator/backend/result/"

var experimentRunCount = 1

// SetExperimentRunCount sets how many times each configuration is executed.
func SetExperimentRunCount(count int) error {
	if count < 1 {
		return fmt.Errorf("experiment run count must be >= 1, got %d", count)
	}
	experimentRunCount = count
	return nil
}

// ExperimentRunCount returns the configured repeat count per experiment setting.
func ExperimentRunCount() int {
	return experimentRunCount
}

type batchSizeSetting struct {
	batchSizeValue   string
	experimentSuffix string
}

var differentBatchSizeSettings = []batchSizeSetting{
	{batchSizeValue: "30", experimentSuffix: "batch_size_90"},
	{batchSizeValue: "100", experimentSuffix: "batch_size_300"},
	{batchSizeValue: "200", experimentSuffix: "batch_size_600"},
	{batchSizeValue: "400", experimentSuffix: "batch_size_1200"},
}

// DifferentBatchSizeConfigurationSettings builds fixed-batch settings with a unique result path prefix.
func DifferentBatchSizeConfigurationSettings(scenarioPrefix, perLinkDelay string) []*entities.ConfigurationSetting {
	settings := make([]*entities.ConfigurationSetting, len(differentBatchSizeSettings))
	for index, batchSizeSetting := range differentBatchSizeSettings {
		settings[index] = &entities.ConfigurationSetting{
			Mapping: map[string]string{
				"per_link_delay":             perLinkDelay,
				"number_of_packets_per_link": batchSizeSetting.batchSizeValue,
				"experiment_name": fmt.Sprintf("%s/%s",
					scenarioPrefix, batchSizeSetting.experimentSuffix),
			},
		}
	}
	return settings
}

// DifferentBatchSizeConfigurationSettingsForDynamicBatch builds dynamic-batch settings with a unique result path prefix.
func DifferentBatchSizeConfigurationSettingsForDynamicBatch(scenarioPrefix, perLinkDelay string) []*entities.ConfigurationSetting {
	settings := make([]*entities.ConfigurationSetting, len(differentBatchSizeSettings))
	for index, batchSizeSetting := range differentBatchSizeSettings {
		settings[index] = &entities.ConfigurationSetting{
			Mapping: map[string]string{
				"per_link_delay":             perLinkDelay,
				"min_ack_for_rtt_estimation": batchSizeSetting.batchSizeValue,
				"experiment_name": fmt.Sprintf("%s/%s",
					scenarioPrefix, batchSizeSetting.experimentSuffix),
			},
		}
	}
	return settings
}

// ConfigurationSettingForRun clones a setting, records run index / seed, and appends run_N when run count > 1.
func ConfigurationSettingForRun(base *entities.ConfigurationSetting, runIndex int) *entities.ConfigurationSetting {
	mapping := make(map[string]string, len(base.Mapping)+2)
	for key, value := range base.Mapping {
		mapping[key] = value
	}
	mapping["experiment_run_index"] = strconv.Itoa(runIndex)
	if experimentRunCount > 1 {
		experimentName := mapping["experiment_name"]
		mapping["experiment_name"] = fmt.Sprintf("%s/run_%d", experimentName, runIndex)
	}
	return &entities.ConfigurationSetting{
		Index:   base.Index,
		Mapping: mapping,
	}
}

func ResultDirForSetting(setting *entities.ConfigurationSetting, experimentIndex int) string {
	if experimentName, ok := setting.Mapping["experiment_name"]; ok {
		return fmt.Sprintf("%s%s", ResultBaseDir, experimentName)
	}
	return fmt.Sprintf("%s%d", ResultBaseDir, experimentIndex)
}

func IsResultExists(resultDir string) bool {
	if !dir.IsDirExists(resultDir) {
		return false
	}
	entries, err := os.ReadDir(resultDir)
	if err != nil {
		return false
	}
	return len(entries) > 0
}
