package main

import (
	"chain_simulation/cmd/root"
	"chain_simulation/cmd/run_normal"
	"chain_simulation/cmd/run_online"
	"os"
)

func main() {
	rootCmd := root.CreateRootCmd()
	runOnlineCmd := run_online.CreateRunOnlineCmd()
	runNormalCmd := run_normal.CreateRunNormalCmd()
	rootCmd.AddCommand(runOnlineCmd)
	rootCmd.AddCommand(runNormalCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
