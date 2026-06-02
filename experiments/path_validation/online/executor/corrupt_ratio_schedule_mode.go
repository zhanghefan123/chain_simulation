package executor

import "fmt"

type CorruptRatioScheduleMode int

const (
	CorruptRatioScheduleRandom CorruptRatioScheduleMode = iota
	CorruptRatioScheduleSequential
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
	default:
		return fmt.Errorf("unknown corrupt-ratio-mode %q, want random or sequential", mode)
	}
	return nil
}
