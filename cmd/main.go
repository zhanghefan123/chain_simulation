package main

import (
	"chain_simulation/cmd/root"
	"chain_simulation/cmd/run_online"
	"chain_simulation/cmd/start"
	"os"
)

func main() {
	rootCmd := root.CreateRootCmd()
	startCmd := start.CreateStartCmd()
	runOnlineCmd := run_online.CreateRunOnlineCmd()
	rootCmd.AddCommand(startCmd, runOnlineCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
