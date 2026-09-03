package executor

import "fmt"

type CorruptRatioScheduleMode int

const (
	CorruptRatioScheduleRandom CorruptRatioScheduleMode = iota
	CorruptRatioScheduleSequential
	CorruptRatioNone
)

var corruptRatioScheduleMode = CorruptRatioScheduleRandom

func SetCorruptRatioScheduleMode(mode CorruptRatioScheduleMode) {
	corruptRatioScheduleMode = mode
}

func GetCorruptRatioScheduleMode() CorruptRatioScheduleMode {
	return corruptRatioScheduleMode
}

func SetCorruptRatioScheduleModeFromString(mode string) error {
	switch mode {
	case "random", "":
		SetCorruptRatioScheduleMode(CorruptRatioScheduleRandom)
	case "sequential", "cyclic":
		SetCorruptRatioScheduleMode(CorruptRatioScheduleSequential)
	case "none":
		SetCorruptRatioScheduleMode(CorruptRatioNone)
	default:
		return fmt.Errorf("unknown corrupt-ratio-mode %q, want random, sequential, or none", mode)
	}
	return nil
}
