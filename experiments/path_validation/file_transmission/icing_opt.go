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

var topologyType = types.TopologyType_SimplePathValidation

func GenerateIcingEvents() ([]*entities.Event, error) {
	currentTime := time.Second * 10
	var icingEvents = []*entities.Event{{
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
	pathValidationProtocols := []string{"ICING", "OPT"}
	for _, pathValidationProtocol := range pathValidationProtocols {
		for currentProcess := 2; currentProcess < processCount; currentProcess += 2 {
			for currentHop := 2; currentHop < hopCount; currentHop += 2 {
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

					// ------------------------------- 添加服务器事件 -------------------------------
					currentTime += time.Second * 300
					serverEvent := &entities.Event{
						StartTime: currentTime,
						Action:    types.ActionType_StartServer,
						Handler: func() error {
							go func() {
								fmt.Printf("current start server %s\n", filePath)
								err = validation_manager.StartServer(serverIndex, currentProcess, pathValidationProtocol, simulationIndex, 31313, "multiprocess_file", networkInterface,
									"IPv4", 1)
								if err != nil {
									fmt.Printf("start server error: %v", err)
								}
							}()
							return nil
						},
					}
					icingEvents = append(icingEvents, serverEvent)
					// ------------------------------- 添加服务器事件 -------------------------------

					// ------------------------------- 添加客户端事件 -------------------------------
					currentTime += time.Second * 1
					clientEvent := &entities.Event{
						StartTime: currentTime,
						Action:    types.ActionType_StartClient,
						Handler: func() error {
							go func() {
								fmt.Printf("current start client %s\n", filePath)
								err = validation_manager.StartClient(1, currentProcess, pathValidationProtocol, 31313, destinations,
									"file", 1024, 512, "",
									0, 0, 0)
								if err != nil {
									fmt.Printf("start client error: %v", err)
								}
							}()
							return nil
						},
					}
					icingEvents = append(icingEvents, clientEvent)
					// ------------------------------- 添加客户端事件 -------------------------------
				}
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
	icingEvents = append(icingEvents, removeEvent)

	currentTime += 40 * time.Second
	waitStopEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_WaitTopologyRemove,
		Handler: func() error {
			return nil
		},
	}
	icingEvents = append(icingEvents, waitStopEvent)
	return icingEvents, nil
}

func IcingOptExperiment() error {
	configurationSettings := []*entities.ConfigurationSetting{
		{
			Mapping: map[string]string{},
		},
	}
	fmt.Printf("Icing Opt experiment\n")
	icingOptEvents, err := GenerateIcingEvents()
	fmt.Printf("number of events: %d\n", len(icingOptEvents))
	if err != nil {
		return fmt.Errorf("error to generate icing/opt events: %v\n", err)
	}

	for _, configurationSetting := range configurationSettings {
		err = experiments.SingleSimulation(configurationSetting, icingOptEvents)
		if err != nil {
			return fmt.Errorf("icing/opt experiment failed: %v", err)
		}
	}

	return nil
}
