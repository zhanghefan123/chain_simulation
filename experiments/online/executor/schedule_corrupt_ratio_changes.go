package executor

import "fmt"

// CorruptRatioScheduleParams controls timed malicious-ratio updates.
type CorruptRatioScheduleParams struct {
	StartTimestamp int
	UpdateInterval int
	MaxUpdateCount int
	LargeRatio     int
	LowRatio       int
}

var (
	Corruptratioscheduleparamsfrequency01s = CorruptRatioScheduleParams{
		StartTimestamp: 5000,
		UpdateInterval: 100,
		MaxUpdateCount: 300,
		LargeRatio:     100000,
		LowRatio:       5000,
	}
	Corruptratioscheduleparamsfrequency05s = CorruptRatioScheduleParams{
		StartTimestamp: 5000,
		UpdateInterval: 500,
		MaxUpdateCount: 60,
		LargeRatio:     100000,
		LowRatio:       5000,
	}
	CorruptRatioScheduleParamsFrequency1s = CorruptRatioScheduleParams{
		StartTimestamp: 5000,
		UpdateInterval: 1000,
		MaxUpdateCount: 30,
		LargeRatio:     100000,
		LowRatio:       5000,
	}
)

// ScheduleCorruptRatioChanges applies the active corrupt-ratio schedule mode.
// CorruptRatioNone skips scheduling and leaves all nodes at their initial ratios.
func ScheduleCorruptRatioChanges(seed int64, params CorruptRatioScheduleParams) error {
	switch GetCorruptRatioScheduleMode() {
	case CorruptRatioScheduleSequential:
		return ScheduleCyclicMaliciousParams(
			params.StartTimestamp,
			params.UpdateInterval,
			params.MaxUpdateCount,
			params.LargeRatio,
			params.LowRatio,
			DefaultCyclicMaliciousNodes,
		)
	case CorruptRatioScheduleRandom:
		return ScheduleRandomMaliciousParams(
			params.StartTimestamp,
			params.UpdateInterval,
			params.MaxUpdateCount,
			params.LargeRatio,
			params.LowRatio,
			seed,
			GetRandomMaliciousBadNodeCount(),
			DefaultRandomMaliciousCandidateNodes,
		)
	case CorruptRatioNone:
		return nil
	default:
		return fmt.Errorf("unsupported corrupt ratio schedule mode: %v", GetCorruptRatioScheduleMode())
	}
}
