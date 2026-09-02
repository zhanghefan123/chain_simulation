package fixed_batch

import (
	"chain_simulation/configs"
	"chain_simulation/entities"
	"chain_simulation/entities/types"
	"chain_simulation/experiments"
	"chain_simulation/experiments/path_validation/online/setup"
	"chain_simulation/modules/breakpoint_awareness"
	"chain_simulation/modules/topology_manager"
	"chain_simulation/modules/validation_manager"
	"chain_simulation/utils/dir"
	"chain_simulation/utils/execute"
	"fmt"
	"strconv"
	"time"
)

var (
	topologyType            = types.TopologyType_SecPathMab
	sourceNodeIndex         = 1      // 这个是写死的
	miniBatchSize           = 30     // 这个是写死的, 而且是 for dynamic batch 的
	learningRate            = 0.2    // 这个是写死的
	minimumDeliveryRatio    = 0.9475 // 这个是写死的
	destinationPort         = 31313  // 这个是写死的
	messageSize             = 512    // 这个是写死的
	interval                = 0.0    // 需要理解这个 interval 是 packet interval 还是 batch interval
	secPathMabType          = types.SecPathMabStrategy_FIXED_BATCH
	enableDadeAlgorithm     = false     // 这个是写死的
	enableDedaAlgorithm     = false     // 这个是写死的
	minAckForRttEstimation  = 100       // 这个是写死的
	experimentTimeElapsedMs = 35 * 1000 // 单位为 ms
)

func GenerateSecPathMabFixedBatchDifferentBatchSizeEvents(setting *entities.ConfigurationSetting) ([]*entities.Event, error) {
	secPathMabBatchEvents := make([]*entities.Event, 0)

	var finalExperimentName = ""
	if experimentName, ok := setting.Mapping["experiment_name"]; ok {
		finalExperimentName = experimentName
	}

	// 获取跳数和每跳的 segment 数
	if numberOfHopsStr, ok := setting.Mapping["number_of_hops"]; ok {
		parseInt, err := strconv.ParseInt(numberOfHopsStr, 10, 64)
		if err != nil {
			return nil, err
		}
		configs.TopConfigInstance.SecPathMabConfig.NumberOfHops = int(parseInt)
	}

	if numberOfSegmentsPerHopStr, ok := setting.Mapping["number_of_segments_per_hop"]; ok {
		parseInt, err := strconv.ParseInt(numberOfSegmentsPerHopStr, 10, 64)
		if err != nil {
			return nil, err
		}
		configs.TopConfigInstance.SecPathMabConfig.NumberOfIntermediateNodes = int(parseInt)
	}

	// 进行重置
	err := setup.UpdateExperimentEnvironment()
	if err != nil {
		return []*entities.Event{}, fmt.Errorf("update experiment environment failed: %v", err)
	}

	// 安装内核模块事件
	currentTime := 5.0 * time.Second
	installKernelEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_InstallKernelModule,
		Handler: func() error {
			go func() {
				err = dir.WithContextManager("/home/zhf/Projects/emulator/path_validation_module", func() error {
					err = execute.Command("sudo", []string{"insmod", "pvm.ko"})
					if err != nil {
						return fmt.Errorf("execute command error due to: %v", err)
					}
					return nil
				})
				if err != nil {
					fmt.Printf("with context manager failed due to: %v", err)
				}
			}()
			return nil
		},
	}
	secPathMabBatchEvents = append(secPathMabBatchEvents, installKernelEvent)

	// 2. 拓扑启动事件
	currentTime += time.Second * 10
	perLinkDelayInMs, err := strconv.ParseFloat(setting.Mapping["per_link_delay"], 64)
	if err != nil {
		return nil, fmt.Errorf("parse per_link_delay failed due to %v", err)
	}
	startTopologyEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_StartTopology,
		Handler: func() error {
			return topology_manager.StartTopology(topologyType,
				&entities.DynamicParameters{
					SecPathMabType: secPathMabType,
					PerLinkDelay:   perLinkDelayInMs,
				})
		},
	}
	secPathMabBatchEvents = append(secPathMabBatchEvents, startTopologyEvent)

	// 3. 进行 osmd 实例的初始化
	var batchSizeInt64 int64
	if batchSizeStr, ok := setting.Mapping["batch_size"]; ok {
		batchSizeInt64, err = strconv.ParseInt(batchSizeStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse number of packets per link error: %w", err)
		}
	}

	var minAckForRttEstimationInt64 int64
	if minAckForRttEstimationInt64Str, ok := setting.Mapping["min_ack_for_rtt_estimation"]; ok {
		minAckForRttEstimationInt64, err = strconv.ParseInt(minAckForRttEstimationInt64Str, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse min ackk for rtt estimation error: %w", err)
		}
		minAckForRttEstimation = int(minAckForRttEstimationInt64)
	}

	var destinations []string
	if destinationStr, ok := setting.Mapping["server_name"]; ok {
		destinations = append(destinations, destinationStr)
	}

	currentTime += 20 * time.Second
	initOsmdEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_InitOsmd,
		Handler: func() error {
			go func() {
				fmt.Printf("current start osmd\n")
				err = validation_manager.InitOsmd(sourceNodeIndex, learningRate, minimumDeliveryRatio, destinationPort, destinations,
					messageSize, int(batchSizeInt64), miniBatchSize, interval, secPathMabType, enableDadeAlgorithm,
					enableDedaAlgorithm, minAckForRttEstimation, experimentTimeElapsedMs)
				if err != nil {
					fmt.Printf("start osmd failed due to: %v", err)
				}
			}()
			return nil
		},
	}
	secPathMabBatchEvents = append(secPathMabBatchEvents, initOsmdEvent)

	// 3. 进行同步
	currentTime += 10 * time.Second
	syncEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_SynchronizeTimestamp,
		Handler: func() error {
			for i := range len(topology_manager.TopologyStartParamsInstance.TopologyParams.Nodes) {
				nodeIndex := i + 1
				go func() {
					err = validation_manager.StartSync(nodeIndex)
					if err != nil {
						fmt.Printf("start synchronize failed due to: %v", err)
					}
				}()
			}
			return nil
		},
	}
	secPathMabBatchEvents = append(secPathMabBatchEvents, syncEvent)

	// 进行 serverIndex 的获取
	var serverIndexInt64 int64
	serverIndex := 0
	if serverIndexStr, ok := setting.Mapping["server_index"]; ok {
		serverIndexInt64, err = strconv.ParseInt(serverIndexStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse server index error: %v", err)
		}
		serverIndex = int(serverIndexInt64)
	}

	// 4. 进行服务器的启动
	currentTime += time.Second * 5
	currentProcess := 1
	pathValidationProtocol := "SEC_PATH_MAB"
	simulationIndex := 1
	networkInterface := fmt.Sprintf("ln%d_idx1", serverIndex)
	serverEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_StartServer,
		Handler: func() error {
			go func() {
				err = validation_manager.StartServer(serverIndex, currentProcess, pathValidationProtocol, simulationIndex,
					31313, "multiprocess_file", networkInterface, "Ipv4",
					1, finalExperimentName)
				if err != nil {
					fmt.Printf("epic start server error: %v", err)
				}
			}()
			return nil
		},
	}
	secPathMabBatchEvents = append(secPathMabBatchEvents, serverEvent)

	// 4. 准备进行 osmd 实例的启动
	currentTime += 5 * time.Second
	startOsmdEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_StartOsmd,
		Handler: func() error {
			go func() {
				err = validation_manager.StartOsmd(sourceNodeIndex)
				if err != nil {
					fmt.Printf("start osmd failed due to: %v", err)
				}
			}()
			return nil
		},
	}
	secPathMabBatchEvents = append(secPathMabBatchEvents, startOsmdEvent)

	// 6. 进行拓扑的清空
	currentTime += 45 * time.Second
	stopEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_StopTopology,
		Handler: func() error {
			return topology_manager.StopTopology()
		},
	}
	secPathMabBatchEvents = append(secPathMabBatchEvents, stopEvent) // 每个节点发送 50 个 , 50 * 3 = 150, 最后一个节点 rtt = 12ms, 0.012 * 25000pkts/s = 300,  大约就是 450 个packet

	// 7. 进行内核模块的卸载
	currentTime += 5 * time.Second
	removeKernelEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_RemoveKernelModule,
		Handler: func() error {
			go func() {
				err = dir.WithContextManager("/home/zhf/Projects/emulator/path_validation_module", func() error {
					err = execute.Command("sudo", []string{"rmmod", "pvm"})
					if err != nil {
						return fmt.Errorf("execute command error due to: %v", err)
					}
					return nil
				})
				if err != nil {
					fmt.Printf("with context manager failed due to: %v", err)
				}
			}()
			return nil
		},
	}
	secPathMabBatchEvents = append(secPathMabBatchEvents, removeKernelEvent)

	// 8. 等待整个拓扑被清空
	currentTime += 20 * time.Second
	waitStopEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_WaitTopologyRemove,
		Handler: func() error {
			return nil
		},
	}
	secPathMabBatchEvents = append(secPathMabBatchEvents, waitStopEvent)
	return secPathMabBatchEvents, nil
}

func GenerateConfigurationSettings() ([]*entities.ConfigurationSetting, error) {
	// 最终的结果
	var configurationSettings []*entities.ConfigurationSetting
	// 选择每次 run 多少次
	runExperimentTimes := 10
	// 选择的 batch size
	batchSizeList := []int{300, 600, 900, 1200}
	// 首先查询自己已经做过的实验
	alreadyRunedExperiments, err := breakpoint_awareness.GetAlreadyCaclulatedFileTransmissionResult(
		"/home/zhf/Projects/emulator/backend/cmd/final_result/")
	if err != nil {
		return nil, fmt.Errorf("get already calculated file transmission result failed: %v", err)
	}
	// 拓扑参数设置 number_of_hops = 2,4,6,8,10, number_of_segments_per_hop = 2
	numberOfSegmentsPerHop := 2
	for currentHop := 2; currentHop <= 10; currentHop += 2 {
		for runExperimentIndex := range runExperimentTimes {
			serverIndex := currentHop*(numberOfSegmentsPerHop+1) + 1
			experimentDestination := fmt.Sprintf("LirNode-%d", serverIndex)
			for _, batchSize := range batchSizeList {
				experimentName := fmt.Sprintf("PROTO:SEC_PATH_MAB|HOP:%d|SEG:%d|BATCH:%d|index:%d", currentHop,
					numberOfSegmentsPerHop, batchSize, runExperimentIndex)
				if _, ok := alreadyRunedExperiments[experimentName]; ok {
					continue
				} else {
					// 如果没有相应的内容的话, 那么才进行 configuration setting 的创建
					configurationSetting := &entities.ConfigurationSetting{
						Mapping: map[string]string{
							"per_link_delay":             "0.0",
							"batch_size":                 strconv.Itoa(batchSize),
							"number_of_segments_per_hop": strconv.Itoa(numberOfSegmentsPerHop),
							"number_of_hops":             strconv.Itoa(currentHop),
							"server_index":               strconv.Itoa(serverIndex),
							"server_name":                experimentDestination,
							"experiment_name":            experimentName,
						},
					}
					configurationSettings = append(configurationSettings, configurationSetting)
				}
			}
		}
	}
	return configurationSettings, nil
}

// RunExperiments 进行多次的实验
func RunExperiments() error {
	configurationSettings, err := GenerateConfigurationSettings()
	if err != nil {
		return fmt.Errorf("generate configuration settings failed: %v", err)
	}
	for _, configurationSetting := range configurationSettings {
		var secPathMabEvents []*entities.Event
		secPathMabEvents, err = GenerateSecPathMabFixedBatchDifferentBatchSizeEvents(configurationSetting)
		err = experiments.SingleSimulation(configurationSetting, secPathMabEvents)
		if err != nil {
			return fmt.Errorf("icing/opt experiment failed: %v", err)
		}
	}
	return nil
}
