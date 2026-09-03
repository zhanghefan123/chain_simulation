package runner

import (
	"chain_simulation/experiments/online/throughput/opt"
	"fmt"
)

func RunAllThroughputExperiments() error {
	runners := []struct {
		name string
		run  func() error
	}{
		//{"fixed_batch/throughput/delay_1_25", fixed_batch.RunExperiments},
		{"opt_experiment", opt.RunExperiments},
	}
	for _, runner := range runners {
		fmt.Printf("start fixed batch throughput scenario: %s\n", runner.name)
		if err := runner.run(); err != nil {
			return fmt.Errorf("fixed batch throughput scenario %s failed: %w", runner.name, err)
		}
		fmt.Printf("finished fixed batch throughput scenario: %s\n", runner.name)
	}
	return nil
}
