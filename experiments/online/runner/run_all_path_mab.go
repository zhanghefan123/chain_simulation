package runner

import (
	pathMabDelay5msFrequency1s "chain_simulation/experiments/online/frequency_1s/path_mab_delay_5ms"
	"fmt"
)

// 	pathMabDelay10msFrequency01s "chain_simulation/experiments/path_validation/online/frequency_0_1s/path_mab_delay_10ms"
//	pathMabDelay125msFrequency01s "chain_simulation/experiments/path_validation/online/frequency_0_1s/path_mab_delay_1_25ms"
//	pathMabDelay25msFrequency01s "chain_simulation/experiments/path_validation/online/frequency_0_1s/path_mab_delay_2_5ms"
//	pathMabDelay5msFrequency01s "chain_simulation/experiments/path_validation/online/frequency_0_1s/path_mab_delay_5ms"
//	pathMabDelay10msFrequency05s "chain_simulation/experiments/path_validation/online/frequency_0_5s/path_mab_delay_10ms"
//	pathMabDelay125msFrequency05s "chain_simulation/experiments/path_validation/online/frequency_0_5s/path_mab_delay_1_25ms"
//	pathMabDelay25msFrequency05s "chain_simulation/experiments/path_validation/online/frequency_0_5s/path_mab_delay_2_5ms"
//	pathMabDelay5msFrequency05s "chain_simulation/experiments/path_validation/online/frequency_0_5s/path_mab_delay_5ms"
//
//	pathMabDelay125msFrequency1s "chain_simulation/experiments/path_validation/online/frequency_1s/path_mab_delay_1_25ms"
//	pathMabDelay25msFrequency1s "chain_simulation/experiments/path_validation/online/frequency_1s/path_mab_delay_2_5ms"
//	pathMabDelay5msFrequency1s "chain_simulation/experiments/path_validation/online/frequency_1s/path_mab_delay_5ms"

func RunAllPathMabExperiments() error {
	runners := []struct {
		name string
		run  func() error
	}{
		//{"path_mab/frequency_0_1s/delay_1_25ms", pathMabDelay125msFrequency01s.SecPathMabPathMabDifferentBatchSizeExperiment},
		//{"path_mab/frequency_0_1s/delay_2_5ms", pathMabDelay25msFrequency01s.SecPathMabPathMabDifferentBatchSizeExperiment},
		//{"path_mab/frequency_0_1s/delay_5ms", pathMabDelay5msFrequency01s.SecPathMabPathMabDifferentBatchSizeExperiment},
		//{"path_mab/frequency_0_1s/delay_10ms", pathMabDelay10msFrequency01s.SecPathMabPathMabDifferentBatchSizeExperiment},

		//{"path_mab/frequency_0_5s/delay_1_25ms", pathMabDelay125msFrequency05s.SecPathMabPathMabDifferentBatchSizeExperiment},
		//{"path_mab/frequency_0_5s/delay_2_5ms", pathMabDelay25msFrequency05s.SecPathMabPathMabDifferentBatchSizeExperiment},
		//{"path_mab/frequency_0_5s/delay_5ms", pathMabDelay5msFrequency05s.SecPathMabPathMabDifferentBatchSizeExperiment},
		//{"path_mab/frequency_0_5s/delay_10ms", pathMabDelay10msFrequency05s.SecPathMabPathMabDifferentBatchSizeExperiment},

		//{"path_mab/frequency_1s/delay_1_25ms", pathMabDelay125msFrequency1s.SecPathMabPathMabDifferentBatchSizeExperiment},
		//{"path_mab/frequency_1s/delay_2_5ms", pathMabDelay25msFrequency1s.SecPathMabPathMabDifferentBatchSizeExperiment},
		{"path_mab/frequency_1s/delay_5ms", pathMabDelay5msFrequency1s.SecPathMabPathMabDifferentBatchSizeExperiment},
		//{"path_mab/frequency_1s/delay_10ms", pathMabDelay10msFrequency1s.SecPathMabPathMabDifferentBatchSizeExperiment},
	}

	for _, runner := range runners {
		fmt.Printf("start path mab scenario: %s\n", runner.name)
		if err := runner.run(); err != nil {
			return fmt.Errorf("path mab scenario %s failed: %w", runner.name, err)
		}
		fmt.Printf("finished path mab scenario: %s\n", runner.name)
	}
	return nil
}
