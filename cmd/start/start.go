package start

import (
	"chain_simulation/configs"
	"chain_simulation/experiments/chainmaker"
	fabrics "chain_simulation/experiments/fabric"
	"chain_simulation/experiments/fiscobcos"
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
			// fabric 实验
			err = fabrics.NormalExperiment()
			if err != nil {
				fmt.Printf("fabric experiment failed: %v", err)
				return
			}
			// chainmaker 实验
			err = chainmaker.WithBlackListExperiment()
			if err != nil {
				fmt.Printf("chainmaker experiment failed: %v", err)
			}
			// fisco bcos 实验
			err = fiscobcos.WithBlackListExperiment()
			if err != nil {
				fmt.Printf("fisco bcos experiment failed: %v", err)
			}
		},
	}
	return createStartCmd
}
