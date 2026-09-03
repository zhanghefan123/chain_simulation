package run_normal

import (
	"chain_simulation/experiments/path_validation/batch_transmission"
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
	err := batch_transmission.FastSelirBatchExperiment()
	if err != nil {
		panic(err)
	}
}
