package runner

import (
	dynamicBatchDelay10msFrequency01s "chain_simulation/experiments/path_validation/online/frequency_0_1s/dynamic_batch_delay_10ms"
	dynamicBatchDelay125msFrequency01s "chain_simulation/experiments/path_validation/online/frequency_0_1s/dynamic_batch_delay_1_25ms"
	dynamicBatchDelay25msFrequency01s "chain_simulation/experiments/path_validation/online/frequency_0_1s/dynamic_batch_delay_2_5ms"
	dynamicBatchDelay5msFrequency01s "chain_simulation/experiments/path_validation/online/frequency_0_1s/dynamic_batch_delay_5ms"
	dynamicBatchDelay10msFrequency05s "chain_simulation/experiments/path_validation/online/frequency_0_5s/dynamic_batch_delay_10ms"
	dynamicBatchDelay125msFrequency05s "chain_simulation/experiments/path_validation/online/frequency_0_5s/dynamic_batch_delay_1_25ms"
	dynamicBatchDelay25msFrequency05s "chain_simulation/experiments/path_validation/online/frequency_0_5s/dynamic_batch_delay_2_5ms"
	dynamicBatchDelay5msFrequency05s "chain_simulation/experiments/path_validation/online/frequency_0_5s/dynamic_batch_delay_5ms"
	dynamicBatchDelay10msFrequency1s "chain_simulation/experiments/path_validation/online/frequency_1s/dynamic_batch_delay_10ms"
	dynamicBatchDelay125msFrequency1s "chain_simulation/experiments/path_validation/online/frequency_1s/dynamic_batch_delay_1_25ms"
	dynamicBatchDelay25msFrequency1s "chain_simulation/experiments/path_validation/online/frequency_1s/dynamic_batch_delay_2_5ms"
	dynamicBatchDelay5msFrequency1s "chain_simulation/experiments/path_validation/online/frequency_1s/dynamic_batch_delay_5ms"
	"fmt"
)

func RunAllDynamicBatchExperiments() error {
	runners := []struct {
		name string
		run  func() error
	}{
		{"dynamic_batch/frequency_0_1s/delay_1_25ms", dynamicBatchDelay125msFrequency01s.SecPathMabDynamicBatchDifferentBatchSizeExperiment},
		{"dynamic_batch/frequency_0_1s/delay_2_5ms", dynamicBatchDelay25msFrequency01s.SecPathMabDynamicBatchDifferentBatchSizeExperiment},
		{"dynamic_batch/frequency_0_1s/delay_5ms", dynamicBatchDelay5msFrequency01s.SecPathMabDynamicBatchDifferentBatchSizeExperiment},
		{"dynamic_batch/frequency_0_1s/delay_10ms", dynamicBatchDelay10msFrequency01s.SecPathMabDynamicBatchDifferentBatchSizeExperiment},
		{"dynamic_batch/frequency_0_5s/delay_1_25ms", dynamicBatchDelay125msFrequency05s.SecPathMabDynamicBatchDifferentBatchSizeExperiment},
		{"dynamic_batch/frequency_0_5s/delay_2_5ms", dynamicBatchDelay25msFrequency05s.SecPathMabDynamicBatchDifferentBatchSizeExperiment},
		{"dynamic_batch/frequency_0_5s/delay_5ms", dynamicBatchDelay5msFrequency05s.SecPathMabDynamicBatchDifferentBatchSizeExperiment},
		{"dynamic_batch/frequency_0_5s/delay_10ms", dynamicBatchDelay10msFrequency05s.SecPathMabDynamicBatchDifferentBatchSizeExperiment},
		{"dynamic_batch/frequency_1s/delay_1_25ms", dynamicBatchDelay125msFrequency1s.SecPathMabDynamicBatchDifferentBatchSizeExperiment},
		{"dynamic_batch/frequency_1s/delay_2_5ms", dynamicBatchDelay25msFrequency1s.SecPathMabDynamicBatchDifferentBatchSizeExperiment},
		{"dynamic_batch/frequency_1s/delay_5ms", dynamicBatchDelay5msFrequency1s.SecPathMabDynamicBatchDifferentBatchSizeExperiment},
		{"dynamic_batch/frequency_1s/delay_10ms", dynamicBatchDelay10msFrequency1s.SecPathMabDynamicBatchDifferentBatchSizeExperiment},
	}

	for _, runner := range runners {
		fmt.Printf("start dynamic batch scenario: %s\n", runner.name)
		if err := runner.run(); err != nil {
			return fmt.Errorf("dynamic batch scenario %s failed: %w", runner.name, err)
		}
		fmt.Printf("finished dynamic batch scenario: %s\n", runner.name)
	}
	return nil
}
