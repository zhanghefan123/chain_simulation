package run_online

import (
	onlineconfig "chain_simulation/experiments/path_validation/online/config"
	onlineexecutor "chain_simulation/experiments/path_validation/online/executor"
	"chain_simulation/experiments/path_validation/online/runner"
	onlinesetup "chain_simulation/experiments/path_validation/online/setup"
	"fmt"

	"github.com/spf13/cobra"
)

func CreateRunOnlineCmd() *cobra.Command {
	var corruptRatioMode string
	var experimentRuns int

	var runOnlineCmd = &cobra.Command{
		Use:   "run-online",
		Short: "run all online path validation experiments",
	}

	var runAllFixedBatchCmd = &cobra.Command{
		Use:   "fixed-batch",
		Short: "run all fixed batch online experiments",
		Run: func(cmd *cobra.Command, args []string) {
			if err := onlineconfig.SetExperimentRunCount(experimentRuns); err != nil {
				fmt.Printf("invalid experiment runs: %v\n", err)
				return
			}
			if err := onlineexecutor.SetCorruptRatioScheduleModeFromString(corruptRatioMode); err != nil {
				fmt.Printf("invalid corrupt-ratio-mode: %v\n", err)
				return
			}
			if err := onlinesetup.InitOnlineExperimentEnvironment(); err != nil {
				fmt.Printf("init online experiment environment failed: %v\n", err)
				return
			}
			if err := runner.RunAllFixedBatchExperiments(); err != nil {
				fmt.Printf("run all fixed batch experiments failed: %v\n", err)
			}
		},
	}
	runAllFixedBatchCmd.Flags().IntVar(&experimentRuns, "runs", 1,
		"number of times to repeat each experiment configuration (results go to .../run_N when > 1)")
	runAllFixedBatchCmd.Flags().StringVar(&corruptRatioMode, "corrupt-ratio-mode", "random",
		"malicious node schedule: random (default) or sequential")

	var runAllDynamicBatchCmd = &cobra.Command{
		Use:   "dynamic-batch",
		Short: "run all dynamic batch online experiments",
		Run: func(cmd *cobra.Command, args []string) {
			if err := onlineconfig.SetExperimentRunCount(experimentRuns); err != nil {
				fmt.Printf("invalid experiment runs: %v\n", err)
				return
			}
			if err := onlineexecutor.SetCorruptRatioScheduleModeFromString(corruptRatioMode); err != nil {
				fmt.Printf("invalid corrupt-ratio-mode: %v\n", err)
				return
			}
			if err := onlinesetup.InitOnlineExperimentEnvironment(); err != nil {
				fmt.Printf("init online experiment environment failed: %v\n", err)
				return
			}
			if err := runner.RunAllDynamicBatchExperiments(); err != nil {
				fmt.Printf("run all dynamic batch experiments failed: %v\n", err)
			}
		},
	}
	runAllDynamicBatchCmd.Flags().IntVar(&experimentRuns, "runs", 1,
		"number of times to repeat each experiment configuration (results go to .../run_N when > 1)")
	runAllDynamicBatchCmd.Flags().StringVar(&corruptRatioMode, "corrupt-ratio-mode", "random",
		"malicious node schedule: random (default) or sequential")

	runOnlineCmd.AddCommand(runAllFixedBatchCmd, runAllDynamicBatchCmd)
	return runOnlineCmd
}
