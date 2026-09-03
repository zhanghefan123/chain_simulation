package run_normal

import (
	"chain_simulation/configs"
	"chain_simulation/experiments/path_validation/batch_transmission"
	"fmt"
	"github.com/spf13/cobra"
)

func CreateRunNormalCmd() *cobra.Command {
	runNormalCmd := &cobra.Command{
		Use:   "run-normal",
		Short: "run all normal experiments",
		Run: func(cmd *cobra.Command, args []string) {
			RunNormalCommand()
			return
		},
	}
	return runNormalCmd
}

func RunNormalCommand() {
	if err := configs.InitTopConfig(); err != nil {
		fmt.Printf("initialize top-level config failed: %v\n", err)
	}
	err := batch_transmission.FastSelirBatchExperiment()
	if err != nil {
		fmt.Printf("run normal experiments failed: %v\n", err)
	}
}
