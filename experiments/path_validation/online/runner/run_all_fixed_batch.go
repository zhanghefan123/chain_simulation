package runner

import (
	fixedBatchDelay5msFrequency1s "chain_simulation/experiments/path_validation/online/frequency_1s/fixed_batch_delay_5ms"
	"fmt"
)

func RunAllFixedBatchExperiments() error {
	runners := []struct {
		name string
		run  func() error
	}{
		//{"fixed_batch/frequency_0_1s/delay_1_25ms", fixedBatchDelay125msFrequency01s.SecPathMabFixedBatchDifferentBatchSizeExperiment},
		//{"fixed_batch/frequency_0_1s/delay_2_5ms", fixedBatchDelay25msFrequency01s.SecPathMabFixedBatchDifferentBatchSizeExperiment},
		//{"fixed_batch/frequency_0_1s/delay_5ms", fixedBatchDelay5msFrequency01s.SecPathMabFixedBatchDifferentBatchSizeExperiment},
		//{"fixed_batch/frequency_0_1s/delay_10ms", fixedBatchDelay10msFrequency01s.SecPathMabFixedBatchDifferentBatchSizeExperiment},
		//{"fixed_batch/frequency_0_5s/delay_1_25ms", fixedBatchDelay125msFrequency05s.SecPathMabFixedBatchDifferentBatchSizeExperiment},
		//{"fixed_batch/frequency_0_5s/delay_2_5ms", fixedBatchDelay25msFrequency05s.SecPathMabFixedBatchDifferentBatchSizeExperiment},
		//{"fixed_batch/frequency_0_5s/delay_5ms", fixedBatchDelay5msFrequency05s.SecPathMabFixedBatchDifferentBatchSizeExperiment},
		//{"fixed_batch/frequency_0_5s/delay_10ms", fixedBatchDelay10msFrequency05s.SecPathMabFixedBatchDifferentBatchSizeExperiment},
		//{"fixed_batch/frequency_1s/delay_1_25ms", fixedBatchDelay125msFrequency1s.SecPathMabFixedBatchDifferentBatchSizeExperiment},
		//{"fixed_batch/frequency_1s/delay_2_5ms", fixedBatchDelay25msFrequency1s.SecPathMabFixedBatchDifferentBatchSizeExperiment},
		{"fixed_batch/frequency_1s/delay_5ms", fixedBatchDelay5msFrequency1s.SecPathMabFixedBatchDifferentBatchSizeExperiment},
		//{"fixed_batch/frequency_1s/delay_10ms", fixedBatchDelay10msFrequency1s.SecPathMabFixedBatchDifferentBatchSizeExperiment},
	}

	for _, runner := range runners {
		fmt.Printf("start fixed batch scenario: %s\n", runner.name)
		if err := runner.run(); err != nil {
			return fmt.Errorf("fixed batch scenario %s failed: %w", runner.name, err)
		}
		fmt.Printf("finished fixed batch scenario: %s\n", runner.name)
	}
	return nil
}
