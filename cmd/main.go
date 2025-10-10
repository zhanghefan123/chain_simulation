package main

import (
	"chain_simulation/cmd/root"
	"chain_simulation/cmd/start"
	"os"
)

func main() {
	rootCmd := root.CreateRootCmd()
	startCmd := start.CreateStartCmd()
	rootCmd.AddCommand(startCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
