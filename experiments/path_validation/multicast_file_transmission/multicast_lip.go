package multicast_file_transmission

import (
	"chain_simulation/entities"
	"chain_simulation/entities/types"
	"chain_simulation/experiments"
	"chain_simulation/modules/breakpoint_awareness"
	"chain_simulation/modules/fast_selir"
	"chain_simulation/modules/topology_manager"
	"chain_simulation/modules/validation_manager"
	"chain_simulation/utils/extract"
	"fmt"
	"time"
)

var topologyType = types.TopologyType_MulticastPathValidation

func GenerateMulticastLiPEvents() ([]*entities.Event, error) {
	currentTime := time.Second * 10
	var multicastLiPEvents = []*entities.Event{{
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
	pathValidationProtocols := []string{"MULTICAST_SELIR"}
	destinationCount := 4

	destinationMapping := map[int][]string{
		1: {"LirNode-3", "LirNode-5"},
		2: {"LirNode-3", "LirNode-5", "LirNode-7"},
		3: {"LirNode-3", "LirNode-5", "LirNode-7", "LirNode-9"},
		4: {"LirNode-3", "LirNode-5", "LirNode-7", "LirNode-9", "LirNode-11"},
	}

	insertedHvfMapping := map[int]int{
		1: 3,
		2: 4,
		3: 5,
		4: 6,
	}

	processCount := 12

	for _, pathValidationProtocol := range pathValidationProtocols {
		for currentDestination := 1; currentDestination <= destinationCount; currentDestination += 1 {
			for currentProcess := 2; currentProcess < processCount; currentProcess += 2 {
				// 结果的记录文件
				filePath := fmt.Sprintf("%s_destinations_%d_processes_%d.txt", pathValidationProtocol, currentDestination, currentProcess)

				// 判断是否已经计算了
				if _, ok := calculatedMapping[filePath]; ok {
					fmt.Println("already calculated")
					continue
				}

				// 进行布隆过滤器大小的修改
				currentTime += time.Second * 10
				modifyBloomFilterEvent := &entities.Event{
					StartTime: currentTime,
					Action:    types.ActionType_ModifyBloomFilter,
					Handler: func() error {
						bfEffectiveBits := fast_selir.CalculateFastSelirBFBits(insertedHvfMapping[currentDestination], 0.00001)
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
				multicastLiPEvents = append(multicastLiPEvents, modifyBloomFilterEvent)

				// 添加服务器事件
				currentTime += time.Second * 10
				serverEvent := &entities.Event{
					StartTime: currentTime,
					Action:    types.ActionType_StartServer,
					Handler: func() error {
						for index, destination := range destinationMapping[currentDestination] {
							if index != 0 {
								go func() {
									fmt.Printf("current start server %s\n", destination)
									destinationIndex, _ := extract.NumberFromString(destination)
									networkInterface := fmt.Sprintf("ln%d_idx1", destinationIndex)
									err = validation_manager.StartServer(destinationIndex, currentProcess, pathValidationProtocol, 0, 31313,
										"multiprocess_file", networkInterface, "IPv4", currentDestination)
									if err != nil {
										fmt.Printf("start server error: %v", err)
									}
								}()
							}
						}
						return nil
					},
				}
				multicastLiPEvents = append(multicastLiPEvents, serverEvent)

				// 添加客户端事件
				currentTime += time.Second * 5
				clientEvent := &entities.Event{
					StartTime: currentTime,
					Action:    types.ActionType_StartClient,
					Handler: func() error {
						go func() {
							fmt.Printf("current start client %s\n", filePath)
							err = validation_manager.StartClient(1, currentProcess, pathValidationProtocol, 31313, destinationMapping[currentDestination],
								"file", 100, 1024, "",
								0, 0, 0)
							if err != nil {
								fmt.Printf("start client error: %v", err)
							}
						}()
						return nil
					},
				}
				multicastLiPEvents = append(multicastLiPEvents, clientEvent)

				// 添加结果处理事件
				// 添加结果处理事件
				currentTime += time.Second * 150
				waitEvent := &entities.Event{
					StartTime: currentTime,
					Action:    types.ActionType_DoNothing,
					Handler: func() error {
						return nil
					},
				}
				multicastLiPEvents = append(multicastLiPEvents, waitEvent)
			}
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
	multicastLiPEvents = append(multicastLiPEvents, removeEvent)

	currentTime += 40 * time.Second
	waitStopEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_WaitTopologyRemove,
		Handler: func() error {
			return nil
		},
	}
	multicastLiPEvents = append(multicastLiPEvents, waitStopEvent)

	return multicastLiPEvents, nil
}

func MulticastLiPExperiment() error {
	configurationSettings := []*entities.ConfigurationSetting{{
		Mapping: map[string]string{},
	}}
	fmt.Printf("multicast LiP experiment")
	multicastLiPBatchEvents, err := GenerateMulticastLiPEvents()
	if err != nil {
		return fmt.Errorf("error to generate multicast LiP batch events: %v\n", err)
	}
	fmt.Printf("number of events: %d\n", len(multicastLiPBatchEvents))
	for _, configurationSetting := range configurationSettings {
		err = experiments.SingleSimulation(configurationSetting, multicastLiPBatchEvents)
		if err != nil {
			return fmt.Errorf("multicast LiP experiment failed: %v", err)
		}
	}
	return nil
}
