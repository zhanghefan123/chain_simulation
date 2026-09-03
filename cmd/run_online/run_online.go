package run_online

import (
	onlineconfig "chain_simulation/experiments/online/config"
	"chain_simulation/experiments/online/executor"
	runner2 "chain_simulation/experiments/online/runner"
	onlinesetup "chain_simulation/experiments/online/setup"
	"fmt"

	"github.com/spf13/cobra"
)

func CreateRunOnlineCmd() *cobra.Command {
	runOnlineCmd := &cobra.Command{
		Use:   "run-online",
		Short: "run all online path validation experiments",
	}

	runOnlineCmd.AddCommand(
		createFixedBatchCmd(),
		createDynamicBatchCmd(),
		createPathMabCmd(),
		createThroughputCmd(),
	)

	return runOnlineCmd
}

func createFixedBatchCmd() *cobra.Command {
	var experimentRuns int
	var corruptRatioMode string
	var badNodeCount int

	cmd := &cobra.Command{
		Use:   "fixed-batch",
		Short: "run all fixed batch online experiments",
		Run: func(cmd *cobra.Command, args []string) {
			if err := prepareOnlineExperiment(experimentRuns, corruptRatioMode, badNodeCount); err != nil {
				fmt.Printf("prepare fixed batch experiment failed: %v\n", err)
				return
			}

			if err := runner2.RunAllFixedBatchExperiments(); err != nil {
				fmt.Printf("run all fixed batch experiments failed: %v\n", err)
			}
		},
	}

	addCommonFlags(cmd, &experimentRuns, &corruptRatioMode, &badNodeCount)

	return cmd
}

func createDynamicBatchCmd() *cobra.Command {
	var experimentRuns int
	var corruptRatioMode string
	var badNodeCount int

	cmd := &cobra.Command{
		Use:   "dynamic-batch",
		Short: "run all dynamic batch online experiments",
		Run: func(cmd *cobra.Command, args []string) {
			if err := prepareOnlineExperiment(experimentRuns, corruptRatioMode, badNodeCount); err != nil {
				fmt.Printf("prepare dynamic batch experiment failed: %v\n", err)
				return
			}

			if err := runner2.RunAllDynamicBatchExperiments(); err != nil {
				fmt.Printf("run all dynamic batch experiments failed: %v\n", err)
			}
		},
	}

	addCommonFlags(cmd, &experimentRuns, &corruptRatioMode, &badNodeCount)

	return cmd
}

func createPathMabCmd() *cobra.Command {
	var experimentRuns int
	var corruptRatioMode string
	var badNodeCount int

	cmd := &cobra.Command{
		Use:   "path-mab",
		Short: "run all path-mab online experiments",
		Run: func(cmd *cobra.Command, args []string) {
			if err := prepareOnlineExperiment(experimentRuns, corruptRatioMode, badNodeCount); err != nil {
				fmt.Printf("prepare path-mab experiment failed: %v\n", err)
				return
			}

			if err := runner2.RunAllPathMabExperiments(); err != nil {
				fmt.Printf("run all path mab experiments failed: %v\n", err)
			}
		},
	}

	addCommonFlags(cmd, &experimentRuns, &corruptRatioMode, &badNodeCount)

	return cmd
}

func createThroughputCmd() *cobra.Command {
	var experimentRuns int
	var corruptRatioMode string

	cmd := &cobra.Command{
		Use:   "throughput",
		Short: "run all throughput online experiments",
		Run: func(cmd *cobra.Command, args []string) {
			if err := prepareThroughputExperiment(experimentRuns, corruptRatioMode); err != nil {
				fmt.Printf("prepare throughput experiment failed: %v\n", err)
				return
			}
			if err := runner2.RunAllThroughputExperiments(); err != nil {
				fmt.Printf("run all throughput experiments failed: %v\n", err)
			}
		},
	}

	cmd.Flags().IntVar(
		&experimentRuns,
		"runs",
		1,
		"number of times to repeat each experiment configuration (results go to .../run_N when > 1)",
	)

	cmd.Flags().StringVar(
		&corruptRatioMode,
		"corrupt-ratio-mode",
		"none",
		"malicious node schedule: none (default for throughput), random, or sequential",
	)

	return cmd
}

func addCommonFlags(
	cmd *cobra.Command,
	experimentRuns *int,
	corruptRatioMode *string,
	badNodeCount *int,
) {
	cmd.Flags().IntVar(
		experimentRuns,
		"runs",
		1,
		"number of times to repeat each experiment configuration (results go to .../run_N when > 1)",
	)

	cmd.Flags().StringVar(
		corruptRatioMode,
		"corrupt-ratio-mode",
		"random",
		"malicious node schedule: random (default), sequential, or none",
	)

	cmd.Flags().IntVar(
		badNodeCount,
		"bad-node-count",
		1,
		"when --corrupt-ratio-mode=random: nodes attacked per update, spread across hops (cannot empty any hop)",
	)
}

func prepareOnlineExperiment(
	experimentRuns int,
	corruptRatioMode string,
	badNodeCount int,
) error {
	if err := onlineconfig.SetExperimentRunCount(experimentRuns); err != nil {
		return fmt.Errorf("invalid experiment runs: %w", err)
	}

	if err := executor.SetCorruptRatioScheduleModeFromString(corruptRatioMode); err != nil {
		return fmt.Errorf("invalid corrupt-ratio-mode: %w", err)
	}

	if err := onlinesetup.InitOnlineExperimentEnvironment(); err != nil {
		return fmt.Errorf("init online experiment environment failed: %w", err)
	}

	if corruptRatioMode == "random" {
		if err := executor.SetRandomMaliciousBadNodeCount(badNodeCount); err != nil {
			return fmt.Errorf("invalid bad-node-count: %w", err)
		}
	}

	return nil
}

func prepareThroughputExperiment(
	experimentRuns int,
	corruptRatioMode string,
) error {
	if err := onlineconfig.SetExperimentRunCount(experimentRuns); err != nil {
		return fmt.Errorf("invalid experiment runs: %w", err)
	}

	if err := executor.SetCorruptRatioScheduleModeFromString(corruptRatioMode); err != nil {
		return fmt.Errorf("invalid corrupt-ratio-mode: %w", err)
	}

	if err := onlinesetup.InitOnlineExperimentEnvironment(); err != nil {
		return fmt.Errorf("init online experiment environment failed: %w", err)
	}

	return nil
}
