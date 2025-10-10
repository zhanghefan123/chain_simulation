package start

import (
	"chain_simulation/configs"
	"chain_simulation/entities/types"
	"chain_simulation/simulation"
	"fmt"

	"github.com/spf13/cobra"
)

func CreateStartCmd() *cobra.Command {
	var createStartCmd = &cobra.Command{
		Use:   "start",
		Short: "start",
		Run: func(cmd *cobra.Command, args []string) {
			err := configs.InitTopConfig()
			if err != nil {
				fmt.Printf("init top config err %v\n", err)
			}
			err = simulation.TickIntervalSimulations(types.TopologyType_ChainMaker)
			if err != nil {
				fmt.Printf("TickIntervalSimulations: %v\n", err)
			}
		},
	}
	return createStartCmd
}
