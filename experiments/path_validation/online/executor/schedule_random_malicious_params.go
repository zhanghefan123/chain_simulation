package executor

import (
	"chain_simulation/entities"
	"chain_simulation/modules/validation_manager"
	"fmt"
	"math/rand"
	"strconv"
)

// DefaultRandomMaliciousCandidateNodes 3 跳 MAB，每跳 2 个分支，共 6 个候选坏节点
var DefaultRandomMaliciousCandidateNodes = []int{2, 3, 5, 6, 8, 9}

const DefaultRandomMaliciousSeed int64 = 1234

// MaliciousSeedForRun derives the RNG seed from the 1-based repeat index of an experiment run.
// run_1 keeps DefaultRandomMaliciousSeed; each additional repeat increments the seed by 1.
func MaliciousSeedForRun(runIndex int) int64 {
	if runIndex < 1 {
		runIndex = 1
	}
	return DefaultRandomMaliciousSeed + int64(runIndex-1)
}

// MaliciousSeedFromSetting reads experiment_run_index from a configuration setting.
func MaliciousSeedFromSetting(setting *entities.ConfigurationSetting) int64 {
	if setting == nil {
		return DefaultRandomMaliciousSeed
	}
	if runIndexStr, ok := setting.Mapping["experiment_run_index"]; ok {
		runIndex, err := strconv.Atoi(runIndexStr)
		if err == nil {
			return MaliciousSeedForRun(runIndex)
		}
	}
	return DefaultRandomMaliciousSeed
}

func ScheduleRandomMaliciousParams(
	startTimestamp, updateInterval, maxUpdateCount, largeRatio, lowRatio int,
	seed int64,
	maliciousCandidateNodes []int,
) error {
	rng := rand.New(rand.NewSource(seed))
	currentTimestamp := startTimestamp + updateInterval
	previousBadNodeIndex := -1
	for currentUpdateCount := 0; currentUpdateCount < maxUpdateCount; currentUpdateCount++ {
		badNodeIndex := maliciousCandidateNodes[rng.Intn(len(maliciousCandidateNodes))]
		for previousBadNodeIndex != -1 && badNodeIndex == previousBadNodeIndex {
			badNodeIndex = maliciousCandidateNodes[rng.Intn(len(maliciousCandidateNodes))]
		}

		if previousBadNodeIndex != -1 {
			err := validation_manager.SetScheduledMaliciousParams(previousBadNodeIndex, currentTimestamp,
				lowRatio, lowRatio, 0, 0)
			if err != nil {
				return fmt.Errorf("reset corrupt ratio failed due to: %v", err)
			}
		}

		err := validation_manager.SetScheduledMaliciousParams(badNodeIndex, currentTimestamp,
			largeRatio, largeRatio, 0, 0)
		if err != nil {
			return fmt.Errorf("change corrupt ratio failed due to: %v", err)
		}
		previousBadNodeIndex = badNodeIndex

		currentTimestamp += updateInterval
	}
	return nil
}
