package test

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestFastSelir(t *testing.T) {
	startTimestamp := 5000
	updateInterval := 1000
	maxUpdateCount := 30
	//largeRatio := 100000
	//lowRatio := 5000
	// 3 跳 MAB，每跳 2 个分支，共 6 个候选坏节点
	maliciousCandidateNodes := []int{2, 3, 5, 6, 8, 9}
	resultList := make([]int, 0)

	rng := rand.New(rand.NewSource(1234))
	currentTimestamp := startTimestamp + updateInterval
	previousBadNodeIndex := -1
	for currentUpdateCount := 0; currentUpdateCount < maxUpdateCount; currentUpdateCount++ {
		badNodeIndex := maliciousCandidateNodes[rng.Intn(len(maliciousCandidateNodes))]
		for previousBadNodeIndex != -1 && badNodeIndex == previousBadNodeIndex {
			badNodeIndex = maliciousCandidateNodes[rng.Intn(len(maliciousCandidateNodes))]
		}

		if previousBadNodeIndex != -1 {
			//err := validation_manager.SetScheduledMaliciousParams(&entities.ScheduledMaliciousParams{
			//	EmployedEpochOrTimestampMs: currentTimestamp,
			//	NodeId: previousBadNodeIndex,
			//	CorruptRatioStart: lowRatio,
			//	CorruptRatioEnd: lowRatio,
			//})
			//if err != nil {
			//	//return fmt.Errorf("reset corrupt ratio failed due to: %v", err)
			//}
		}

		//err := validation_manager.SetScheduledMaliciousParams(&entities.ScheduledMaliciousParams{
		//	EmployedEpochOrTimestampMs: currentTimestamp,
		//	NodeId: badNodeIndex,
		//	CorruptRatioStart: largeRatio,
		//	CorruptRatioEnd: largeRatio,
		//})
		//if err != nil {
		//return fmt.Errorf("change corrupt ratio failed due to: %v", err)
		//}
		resultList = append(resultList, badNodeIndex)
		previousBadNodeIndex = badNodeIndex

		currentTimestamp += updateInterval
	}
	fmt.Println(resultList)
}
