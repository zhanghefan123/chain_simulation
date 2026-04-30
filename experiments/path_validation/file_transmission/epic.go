package file_transmission

import (
	"chain_simulation/entities"
	"chain_simulation/entities/types"
	"chain_simulation/experiments"
	"chain_simulation/modules/breakpoint_awareness"
	"chain_simulation/modules/topology_manager"
	"chain_simulation/modules/validation_manager"
	"fmt"
	"time"
)

// GenerateEpicEvents 进行 Epic Events 的生成
func GenerateEpicEvents() ([]*entities.Event, error) {
	currentTime := time.Second * 10
	var epicEvents = []*entities.Event{{
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
	pathValidationProtocol := "Epic"
	simulationCount := 1
	processCount := 12
	hopCount := 11
	// index 等于 4, 6, 8, 10, 12
	for currentProcess := 2; currentProcess < processCount; currentProcess += 2 {
		for currentHop := 2; currentHop < hopCount; currentHop += 2 {
			for simulationIndex := 0; simulationIndex < simulationCount; simulationIndex++ {
				serverIndex := 1 + currentHop
				networkInterface := fmt.Sprintf("ln%d_idx1", serverIndex)
				serverName := fmt.Sprintf("LirNode-%d", serverIndex)
				sessionSetupDestinations := []string{"LirNode-1"}
				destinations := []string{serverName}

				// 获取在这个参数配置下的文件名称, 如果已经计算了就直接进行返回
				filePath := fmt.Sprintf("%s_result_node_name_%s_processes_%d_index_%d.txt", pathValidationProtocol, serverName, currentProcess, simulationIndex)
				if _, ok := calculatedMapping[filePath]; ok {
					fmt.Println("already calculated continue")
					continue
				}

				// ------------------------------- 添加会话建立报文 -------------------------------
				currentTime += time.Second * 200
				sessionSetupEvent := &entities.Event{
					StartTime: currentTime,
					Action:    types.ActionType_StartClient,
					Handler: func() error {
						go func() {
							// content 如果为 "" 则会启用 input 要求用户进行输入
							err = validation_manager.StartClient(serverIndex, currentProcess, "EPIC_SESSION_SETUP", 31313,
								sessionSetupDestinations, "single", 1024, 1024, "a",
								0, 0, 0)
							if err != nil {
								fmt.Printf("epic start client error: %v", err)
							}
						}()
						return nil
					},
				}
				epicEvents = append(epicEvents, sessionSetupEvent)
				// ------------------------------- 添加会话建立报文 -------------------------------

				// ------------------------------- 添加服务器事件 -------------------------------
				currentTime += time.Second * 5
				serverEvent := &entities.Event{
					StartTime: currentTime,
					Action:    types.ActionType_StartServer,
					Handler: func() error {
						go func() {
							err = validation_manager.StartServer(serverIndex, currentProcess, pathValidationProtocol, simulationIndex, 31313, "multiprocess_file", networkInterface, "IPv4", 1)
							if err != nil {
								fmt.Printf("epic start server error: %v", err)
							}
						}()
						return nil
					},
				}
				epicEvents = append(epicEvents, serverEvent)
				// ------------------------------- 添加服务器事件 -------------------------------

				// ------------------------------- 添加客户端事件 -------------------------------
				currentTime += time.Second * 1
				clientEvent := &entities.Event{
					StartTime: currentTime,
					Action:    types.ActionType_StartClient,
					Handler: func() error {
						go func() {
							// 对于 file 传输模式, 没有关系
							err = validation_manager.StartClient(1, currentProcess, "EPIC_DATA", 31313, destinations,
								"file", 1024, 512, "",
								0, 0, 0)
							if err != nil {
								fmt.Printf("epic start client error: %v", err)
							}
						}()
						return nil
					},
				}
				epicEvents = append(epicEvents, clientEvent)
				// ------------------------------- 添加客户端事件 -------------------------------
			}
		}
	}

	currentTime += 200 * time.Second
	removeEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_StopTopology,
		Handler: func() error {
			return topology_manager.StopTopology()
		},
	}
	epicEvents = append(epicEvents, removeEvent)

	currentTime += 40 * time.Second
	waitStopEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_WaitTopologyRemove,
		Handler: func() error {
			return nil
		},
	}
	epicEvents = append(epicEvents, waitStopEvent)
	return epicEvents, nil
}

func EpicExperiment() error {
	configurationSettings := []*entities.ConfigurationSetting{
		{
			Mapping: map[string]string{},
		},
	}

	fmt.Printf("Epic experiment\n")
	epicEvents, err := GenerateEpicEvents()
	fmt.Printf("number of events: %d\n", len(epicEvents))
	if err != nil {
		fmt.Printf("error to generate epic events: %v\n", err)
	}

	for _, configurationSetting := range configurationSettings {
		err = experiments.SingleSimulation(configurationSetting, epicEvents)
		if err != nil {
			return fmt.Errorf("epic experiment failed: %v", err)
		}
	}

	return nil
}
