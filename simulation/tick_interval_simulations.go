package simulation

import (
	"chain_simulation/entities/types"
	"fmt"
)

func TickIntervalSimulations(topologyType types.TopologyType) error {
	key := "tick_interval_ms"
	values := []string{"1", "5", "10", "20"}
	for index, value := range values {
		err := SingleSimulation(key, value, topologyType, index+1)
		if err != nil {
			return fmt.Errorf("tick interval simulation error: %v", err)
		}
	}
	fmt.Printf("Tick interval simulations finished\n")
	return nil
}
