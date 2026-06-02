package executor

import (
	"chain_simulation/modules/validation_manager"
	"fmt"
)

// DefaultCyclicMaliciousNodes 按 hop 顺序依次破坏: 2 -> 5 -> 8 -> 3 -> 6 -> 9
var DefaultCyclicMaliciousNodes = []int{2, 5, 8, 3, 6, 9}

func ScheduleCyclicMaliciousParams(startTimestamp, updateInterval, maxUpdateCount, largeRatio, lowRatio int, maliciousNodes []int) error {
	currentTimestamp := startTimestamp + updateInterval
	previousBadNodeIndex := -1
	for currentUpdateCount := 0; currentUpdateCount < maxUpdateCount; currentUpdateCount++ {
		badNodeIndex := maliciousNodes[currentUpdateCount%len(maliciousNodes)]

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
