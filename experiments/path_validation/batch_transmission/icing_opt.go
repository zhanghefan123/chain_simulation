package batch_transmission

import (
	"chain_simulation/entities"
	"chain_simulation/entities/types"
	"chain_simulation/experiments"
	"chain_simulation/modules/breakpoint_awareness"
	"chain_simulation/modules/topology_manager"
	"chain_simulation/modules/validation_manager"
	"chain_simulation/utils/file"
	"fmt"
	"time"
)

var topologyType = types.TopologyType_SimplePathValidation

func GenerateIcingOptBatchEvents() ([]*entities.Event, error) {
	currentTime := time.Second * 10
	var icingOptBatchEvents = []*entities.Event{{
		StartTime: currentTime,
		Action:    types.ActionType_StartTopology,
		Handler: func() error {
			return topology_manager.StartTopology(topologyType, &entities.DynamicParameters{ConsensusThreadCount: 0})
		},
	}}
	currentTime += time.Second * 20
	clearKernelLogEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_ClearKernelLog,
		Handler: func() error {
			go func() {
				fmt.Printf("clear kernel log\n")
				kernelLogFilePath := "/var/log/kern.log"
				// 进行原始文件的清空
				err := file.ClearFile(kernelLogFilePath)
				if err != nil {
					fmt.Printf("cannot clear file: %v", err)
				}
			}()
			return nil
		},
	}
	icingOptBatchEvents = append(icingOptBatchEvents, clearKernelLogEvent)

	calculateMapping, err := breakpoint_awareness.GetAlreadyCaclulatedFileTransmissionResult("/home/zhf/Projects/emulator/backend/cmd/final_result")
	if err != nil {
		return nil, fmt.Errorf("get break point failed due to: %v", err)
	}
	pathValidationProtocols := []string{"OPT"}
	hopCount := 11
	for _, pathValidationProtocol := range pathValidationProtocols {
		for currentHop := 2; currentHop < hopCount; currentHop += 2 {
			// 获取服务器索引
			serverIndex := 1 + currentHop
			// 获取服务器名称
			serverName := fmt.Sprintf("LirNode-%d", serverIndex)
			// 获取接口名称
			networkInterface := fmt.Sprintf("ln%d_idx1", serverIndex)
			// 获取在这个参数下的文件名称
			filePath := fmt.Sprintf("%s_batch_result_node_name_%s.txt", pathValidationProtocol, serverName)
			// 判断是否已经计算过了
			if _, ok := calculateMapping[filePath]; ok {
				fmt.Println("already calculated")
				continue
			}
			destinations := []string{serverName}
			// ----------------- 添加服务器事件 -----------------
			currentTime += time.Second * 10
			serverEvent := &entities.Event{
				StartTime: currentTime,
				Action:    types.ActionType_StartServer,
				Handler: func() error {
					go func() {
						fmt.Printf("current start server %s\n", filePath)
						err = validation_manager.StartServer(serverIndex, 1, pathValidationProtocol, 0, 31313,
							"text", networkInterface, "IPv4", 1)
						if err != nil {
							fmt.Printf("start server error: %v", err)
						}
					}()
					return nil
				},
			}
			icingOptBatchEvents = append(icingOptBatchEvents, serverEvent)
			// ----------------- 添加服务器事件 -----------------

			// ----------------- 添加客户端事件 -----------------
			currentTime += time.Second * 1
			clientEvent := &entities.Event{
				StartTime: currentTime,
				Action:    types.ActionType_StartClient,
				Handler: func() error {
					go func() {
						batchSize := 1000
						messageSize := 1024
						interval := 0.1
						fmt.Printf("current start client %s\n", filePath)
						err = validation_manager.StartClient(1, 1, pathValidationProtocol, 31313, destinations,
							"batch", 1024, 1024, "",
							batchSize, messageSize, interval)
						if err != nil {
							fmt.Printf("start client error: %v", err)
						}
					}()
					return nil
				},
			}
			icingOptBatchEvents = append(icingOptBatchEvents, clientEvent)
			// ----------------- 添加客户端事件 -----------------

			// ----------------- 添加结果处理事件 ---------------------
			currentTime += time.Second * 150
			resultProcessingEvent := &entities.Event{
				StartTime: currentTime,
				Action:    types.ActionType_ResultHandling,
				Handler: func() error {
					kernelLogFilePath := "/var/log/kern.log"
					err = file.CopyFileWithName(kernelLogFilePath, "/home/zhf/Projects/emulator/backend/cmd/final_result/", filePath)
					if err != nil {
						fmt.Printf("copy file with name error: %v", err)
					}
					// 进行原始文件的清空
					err = file.ClearFile(kernelLogFilePath)
					if err != nil {
						fmt.Printf("cannot clear file: %v", err)
					}
					return nil
				},
			}
			icingOptBatchEvents = append(icingOptBatchEvents, resultProcessingEvent)
			// ----------------- 添加结果处理事件 ---------------------
		}
	}

	currentTime += 40 * time.Second
	removeEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_StopTopology,
		Handler: func() error {
			return topology_manager.StopTopology()
		},
	}
	icingOptBatchEvents = append(icingOptBatchEvents, removeEvent)

	currentTime += 40 * time.Second
	waitStopEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_WaitTopologyRemove,
		Handler: func() error {
			return nil
		},
	}
	icingOptBatchEvents = append(icingOptBatchEvents, waitStopEvent)
	return icingOptBatchEvents, nil
}

func IcingOptBatchExperiment() error {
	configurationSettings := []*entities.ConfigurationSetting{{
		Mapping: map[string]string{},
	}}
	fmt.Printf("icing opt experiment")
	icingOptBatchEvents, err := GenerateIcingOptBatchEvents()
	if err != nil {
		return fmt.Errorf("error to generate icing/opt batch events: %v\n", err)
	}
	fmt.Printf("number of events: %d\n", len(icingOptBatchEvents))
	for _, configurationSetting := range configurationSettings {
		err = experiments.SingleSimulation(configurationSetting, icingOptBatchEvents)
		if err != nil {
			return fmt.Errorf("icing/opt experiment failed: %v", err)
		}
	}
	return nil
}
