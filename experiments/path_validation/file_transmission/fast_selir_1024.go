package file_transmission

import (
	"chain_simulation/entities"
	"chain_simulation/entities/types"
	"chain_simulation/experiments"
	"chain_simulation/modules/breakpoint_awareness"
	"chain_simulation/modules/fast_selir"
	"chain_simulation/modules/topology_manager"
	"chain_simulation/modules/validation_manager"
	"fmt"
	"time"
)

// GenerateFastSelir1024Events 生成 fast_selir_events 的实验
func GenerateFastSelir1024Events() ([]*entities.Event, error) {
	currentTime := time.Second * 10
	var fastSelirEvents = []*entities.Event{{
		StartTime: currentTime,
		Action:    types.ActionType_StartTopology,
		Handler: func() error {
			return topology_manager.StartTopology(topologyType, &entities.DynamicParameters{ConsensusThreadCount: 0})
		},
	}}
	calculatedMapping, err := breakpoint_awareness.GetAlreadyCaclulatedFileTransmissionResult("/home/zhf/Projects/emulator/backend/cmd/final_result")
	if err != nil {
		return nil, fmt.Errorf("get breakpoint failed: %v", err)
	}
	simulationCount := 1
	processCount := 12
	hopCount := 11
	pathValidationProtocols := []string{"FAST_SELIR"}
	for _, pathValidationProtocol := range pathValidationProtocols {
		for currentHop := 2; currentHop < hopCount; currentHop += 2 {
			for currentProcess := 2; currentProcess < processCount; currentProcess += 2 {
				for simulationIndex := 0; simulationIndex < simulationCount; simulationIndex++ {
					serverIndex := 1 + currentHop
					networkInterface := fmt.Sprintf("ln%d_idx1", serverIndex)
					serverName := fmt.Sprintf("LirNode-%d", serverIndex)
					destinations := []string{serverName}
					// 获取在这个参数配置下的文件名称
					filePath := fmt.Sprintf("%s_result_node_name_%s_processes_%d_index_%d.txt", pathValidationProtocol, serverName, currentProcess, simulationIndex)
					if _, ok := calculatedMapping[filePath]; ok {
						fmt.Println("already calculated continue")
						continue
					}

					// ------------------------------- 进行所有的布隆过滤器的配置 -------------------------------
					currentTime += time.Second * 200
					modifyBloomFilterEvent := &entities.Event{
						StartTime: currentTime,
						Action:    types.ActionType_ModifyBloomFilter,
						Handler: func() error {
							bfEffectiveBits := fast_selir.CalculateFastSelirBFBits(currentHop, 0.00001)
							go func() {
								for index := 1; index <= 1; index++ {
									err = validation_manager.ModifyBloomFilter(index, bfEffectiveBits)
									if err != nil {
										fmt.Printf("modify bloom filter failed: %v", err)
									}
								}
							}()
							return nil
						},
					}
					fastSelirEvents = append(fastSelirEvents, modifyBloomFilterEvent)
					// ------------------------------- 添加服务器事件 -------------------------------

					currentTime += time.Second * 25
					serverEvent := &entities.Event{
						StartTime: currentTime,
						Action:    types.ActionType_StartServer,
						Handler: func() error {
							go func() {
								fmt.Printf("current start server %s\n", filePath)
								err = validation_manager.StartServer(serverIndex, currentProcess, pathValidationProtocol, simulationIndex, 31313, "multiprocess_file", networkInterface, "IPv4", 1)
								if err != nil {
									fmt.Printf("start server error: %v", err)
								}
							}()
							return nil
						},
					}
					fastSelirEvents = append(fastSelirEvents, serverEvent)
					// ------------------------------- 添加服务器事件 -------------------------------

					// ------------------------------- 添加客户端事件 -------------------------------
					currentTime += time.Second * 1
					clientEvent := &entities.Event{
						StartTime: currentTime,
						Action:    types.ActionType_StartClient,
						Handler: func() error {
							go func() {
								fmt.Printf("current start client %s\n", filePath)
								err = validation_manager.StartClient(1, currentProcess, pathValidationProtocol, 31313,
									destinations, "file", 1024, 1024, "",
									0, 0, 0)
								if err != nil {
									fmt.Printf("start client error: %v", err)
								}
							}()
							return nil
						},
					}
					fastSelirEvents = append(fastSelirEvents, clientEvent)
					// ------------------------------- 添加客户端事件 -------------------------------
				}
			}
		}
	}

	// ------------------------------- 进行后端的删除  -------------------------------
	currentTime += time.Second * 200
	backendStopEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_StopTopology,
		Handler: func() error {
			return topology_manager.StopTopology()
		},
	}
	fastSelirEvents = append(fastSelirEvents, backendStopEvent)
	// ------------------------------- 进行后端的删除  -------------------------------

	// ------------------------------- 等待后端删除完成 ------------------------------
	currentTime += time.Second * 40
	waitStopEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_WaitTopologyRemove,
		Handler: func() error {
			return nil
		},
	}
	fastSelirEvents = append(fastSelirEvents, waitStopEvent)
	// ------------------------------- 等待后端删除完成 ------------------------------

	return fastSelirEvents, nil
}

func FastSelir1024Experiment() error {
	configurationSettings := []*entities.ConfigurationSetting{
		{
			Mapping: map[string]string{},
		},
	}

	fmt.Printf("fast selir 1024 experiment\n")
	fastSelirEvents, err := GenerateFastSelir1024Events()
	fmt.Printf("number of events: %d\n", len(fastSelirEvents))
	if err != nil {
		fmt.Printf("error to generate fast selir events: %v\n", err)
	}

	for _, configurationSetting := range configurationSettings {
		err = experiments.SingleSimulation(configurationSetting, fastSelirEvents)
		if err != nil {
			return fmt.Errorf("fast selir experiment failed: %v", err)
		}
	}

	return nil
}
