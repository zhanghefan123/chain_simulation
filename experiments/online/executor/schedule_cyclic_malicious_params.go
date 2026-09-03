package executor

import (
	"chain_simulation/entities"
	"chain_simulation/modules/validation_manager"
	"fmt"
)

func ScheduleCyclicMaliciousParams(startTimestamp, updateInterval, maxUpdateCount, largeRatio, lowRatio int, maliciousNodes []int) error {
	// Ensure candidates reflect current topology settings when possible.
	RefreshMaliciousCandidatesFromConfig()
	if len(maliciousNodes) == 0 {
		maliciousNodes = DefaultCyclicMaliciousNodes
	}

	currentTimestamp := startTimestamp + updateInterval
	previousBadNodeIndex := -1
	for currentUpdateCount := 0; currentUpdateCount < maxUpdateCount; currentUpdateCount++ {
		badNodeIndex := maliciousNodes[currentUpdateCount%len(maliciousNodes)]

		if previousBadNodeIndex != -1 {
			err := validation_manager.SetScheduledMaliciousParams(&entities.ScheduledMaliciousParams{
				EmployedEpochOrTimestampMs:     currentTimestamp,
				NodeId:                         previousBadNodeIndex,
				CorruptRatioStart:              lowRatio,
				CorruptRatioEnd:                lowRatio,
				CorruptSpecialPacketRatioStart: 0,
				CorruptSpecialPacketRatioEnd:   0,
			})

			if err != nil {
				return fmt.Errorf("reset corrupt ratio failed due to: %v", err)
			}
		}

		err := validation_manager.SetScheduledMaliciousParams(&entities.ScheduledMaliciousParams{
			EmployedEpochOrTimestampMs:     currentTimestamp,
			NodeId:                         badNodeIndex,
			CorruptRatioStart:              largeRatio,
			CorruptRatioEnd:                largeRatio,
			CorruptSpecialPacketRatioStart: 0,
			CorruptSpecialPacketRatioEnd:   0,
		})

		if err != nil {
			return fmt.Errorf("change corrupt ratio failed due to: %v", err)
		}
		previousBadNodeIndex = badNodeIndex

		currentTimestamp += updateInterval
	}
	return nil
}
