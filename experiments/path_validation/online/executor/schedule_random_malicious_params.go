package executor

import (
	"chain_simulation/entities"
	"chain_simulation/modules/validation_manager"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
)

type maliciousParamOp struct {
	nodeIndex int
	timestamp int
	ratio     int
	reset     bool // only affects error message text
}

const DefaultRandomMaliciousSeed int64 = 1234

// 进行恶意的 seed 的设置
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
	badNodeCount int,
	maliciousCandidateNodes []int,
) error {
	// Ensure candidates reflect current topology settings when possible.
	RefreshMaliciousCandidatesFromConfig()

	if err := ValidateRandomMaliciousBadNodeCount(badNodeCount, DefaultRandomMaliciousCandidateGroups); err != nil {
		return err
	}
	_ = maliciousCandidateNodes // keep parameter for backward compatibility/experiments

	rng := rand.New(rand.NewSource(seed))
	currentTimestamp := startTimestamp + updateInterval
	previousBadNodes := make([]int, 0, badNodeCount)
	ops := make([]maliciousParamOp, 0, maxUpdateCount*badNodeCount*2)
	for currentUpdateCount := 0; currentUpdateCount < maxUpdateCount; currentUpdateCount++ {
		// 1) spread attacks across hop groups without attacking every node in any group
		badNodes := pickRandomBadNodes(rng, DefaultRandomMaliciousCandidateGroups, badNodeCount)

		// 2) avoid selecting exactly the same set consecutively
		if len(previousBadNodes) > 0 && len(previousBadNodes) == len(badNodes) {
			same := true
			for i := range badNodes {
				if badNodes[i] != previousBadNodes[i] {
					same = false
					break
				}
			}
			if same {
				badNodes = pickRandomBadNodes(rng, DefaultRandomMaliciousCandidateGroups, badNodeCount)
			}
		}

		// 3) reset previous bad nodes back to low ratio (skip nodes also selected this round;
		//    serial order was reset-then-set, so only the large-ratio schedule is needed)
		badSet := make(map[int]struct{}, len(badNodes))
		for _, bad := range badNodes {
			badSet[bad] = struct{}{}
		}
		for _, prev := range previousBadNodes {
			if _, stillBad := badSet[prev]; stillBad {
				continue
			}
			ops = append(ops, maliciousParamOp{
				nodeIndex: prev,
				timestamp: currentTimestamp,
				ratio:     lowRatio,
				reset:     true,
			})
		}

		// 4) set new bad nodes to large ratio
		for _, bad := range badNodes {
			ops = append(ops, maliciousParamOp{
				nodeIndex: bad,
				timestamp: currentTimestamp,
				ratio:     largeRatio,
			})
		}
		previousBadNodes = append(previousBadNodes[:0], badNodes...)

		currentTimestamp += updateInterval
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var scheduleErr error
	for _, op := range ops {
		wg.Add(1)
		go func(op maliciousParamOp) {
			defer wg.Done()
			err := validation_manager.SetScheduledMaliciousParams(op.nodeIndex, op.timestamp,
				op.ratio, op.ratio, 0, 0)
			if err == nil {
				return
			}
			mu.Lock()
			if scheduleErr == nil {
				if op.reset {
					scheduleErr = fmt.Errorf("reset corrupt ratio failed due to: %v", err)
				} else {
					scheduleErr = fmt.Errorf("change corrupt ratio failed due to: %v", err)
				}
			}
			mu.Unlock()
		}(op)
	}
	wg.Wait()
	return scheduleErr
}
