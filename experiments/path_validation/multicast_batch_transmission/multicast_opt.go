package multicast_batch_transmission

import (
	"chain_simulation/entities"
	"chain_simulation/entities/types"
	"chain_simulation/experiments"
	"chain_simulation/modules/experiment_related/breakpoint_awareness"
	"chain_simulation/modules/topology_manager"
	"chain_simulation/modules/validation_manager"
	"chain_simulation/utils/file"
	"fmt"
	"time"
)

var topologyType = types.TopologyType_MulticastPathValidation

func GenerateMulticastOptBatchEvents() ([]*entities.Event, error) {
	currentTime := time.Second * 10
	var multicastOptBatchEvents = []*entities.Event{{
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
	multicastOptBatchEvents = append(multicastOptBatchEvents, clearKernelLogEvent)
	calculateMapping, err := breakpoint_awareness.GetAlreadyCaclulatedResult("/home/zhf/Projects/emulator/backend/cmd/final_result")
	if err != nil {
		return nil, fmt.Errorf("get break point failed due to: %v", err)
	}
	pathValidationProtocols := []string{"MULTICAST_OPT"}
	destinationCount := 4

	destinationMapping := map[int][]string{
		1: {"LirNode-3", "LirNode-5"},
		2: {"LirNode-3", "LirNode-5", "LirNode-7"},
		3: {"LirNode-3", "LirNode-5", "LirNode-7", "LirNode-9"},
		4: {"LirNode-3", "LirNode-5", "LirNode-7", "LirNode-9", "LirNode-11"},
	}
	for _, pathValidationProtocol := range pathValidationProtocols {
		for currentDestination := 1; currentDestination <= destinationCount; currentDestination += 1 {
			// 结果的记录文件
			filePath := fmt.Sprintf("%s_%d_destinations_batch_result.txt", pathValidationProtocol, currentDestination)

			// 判断是否已经计算了
			if _, ok := calculateMapping[filePath]; ok {
				fmt.Println("already calculated")
				continue
			}

			// 添加客户端事件
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
						err = validation_manager.StartClient(1, &entities.StartClient{
							SelectedNetworkLayer: pathValidationProtocol,
							DestinationPort:      31313,
							Processes:            1,
							Destinations:         destinationMapping[currentDestination],
							TransmissionPattern:  "batch",
							FileSize:             1024,
							BufferSize:           1024,
							MessageSize:          messageSize,
							BatchSize:            batchSize,
							Interval:             interval,
						})
						if err != nil {
							fmt.Printf("start client error: %v", err)
						}
					}()
					return nil
				},
			}
			multicastOptBatchEvents = append(multicastOptBatchEvents, clientEvent)

			// 添加结果处理事件
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
			multicastOptBatchEvents = append(multicastOptBatchEvents, resultProcessingEvent)
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
	multicastOptBatchEvents = append(multicastOptBatchEvents, removeEvent)

	currentTime += 40 * time.Second
	waitStopEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_WaitTopologyRemove,
		Handler: func() error {
			return nil
		},
	}
	multicastOptBatchEvents = append(multicastOptBatchEvents, waitStopEvent)

	return multicastOptBatchEvents, nil
}

func MulticastOptBatchExperiment() error {
	configurationSettings := []*entities.ConfigurationSetting{{
		Mapping: map[string]string{},
	}}
	fmt.Printf("multicast opt experiment")
	multicastOptBatchEvents, err := GenerateMulticastOptBatchEvents()
	if err != nil {
		return fmt.Errorf("error to generate multicast opt batch events: %v\n", err)
	}
	fmt.Printf("number of events: %d\n", len(multicastOptBatchEvents))
	for _, configurationSetting := range configurationSettings {
		err = experiments.SingleSimulation(configurationSetting, multicastOptBatchEvents)
		if err != nil {
			return fmt.Errorf("multicast opt experiment failed: %v", err)
		}
	}
	return nil
}
