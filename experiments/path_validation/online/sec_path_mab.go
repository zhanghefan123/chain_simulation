package online

import (
	"chain_simulation/entities"
	"chain_simulation/entities/types"
	"chain_simulation/experiments"
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
	sourceNodeIndex         = 1
	numberOfEpochs          = 500
	numberOfPacketsPerLink  = 100
	miniBatchSize           = 100
	learningRate            = 0.2
	minimumDeliveryRatio    = 0.8
	destinationPort         = 31313
	destinations            = []string{"LirNode-10"}
	messageSize             = 512
	interval                = 0.00005 // 需要理解这个 interval 是 packet interval 还是 batch interval
	secPathMabType          = types.SecPathMabStrategy_FIXED_BATCH
	enableDadeAlgorithm     = false
	enableDedaAlgorithm     = false
	minAckForRttEstimation  = 100
	experimentTimeElapsedMs = 30 * 1000 // 单位为 ms
)

func GenerateSecPathMabEvents(currentExperimentIndex int, setting *entities.ConfigurationSetting) ([]*entities.Event, error) {
	secPathMabBatchEvents := make([]*entities.Event, 0)

	// 安装内核模块事件
	currentTime := 5.0 * time.Second
	installKernelEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_InstallKernelModule,
		Handler: func() error {
			go func() {
				err := dir.WithContextManager("/home/zhf/Projects/emulator/path_validation_module", func() error {
					err := execute.Command("sudo", []string{"insmod", "pvm.ko"})
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
	var numberOfPacketsPerLinkInt64 int64
	if numberOfPacketsPerLinkStr, ok := setting.Mapping["number_of_packets_per_link"]; ok {
		numberOfPacketsPerLinkInt64, err = strconv.ParseInt(numberOfPacketsPerLinkStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse number of packets per link error: %w", err)
		}
		numberOfPacketsPerLink = int(numberOfPacketsPerLinkInt64)
	}

	var minAckForRttEstimationInt64 int64
	if minAckForRttEstimationInt64Str, ok := setting.Mapping["min_ack_for_rtt_estimation"]; ok {
		minAckForRttEstimationInt64, err = strconv.ParseInt(minAckForRttEstimationInt64Str, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse min ackk for rtt estimation error: %w", err)
		}
		minAckForRttEstimation = int(minAckForRttEstimationInt64)
	}

	currentTime += 20 * time.Second
	initOsmdEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_InitOsmd,
		Handler: func() error {
			go func() {
				fmt.Printf("current start osmd\n")

				err = validation_manager.InitOsmd(sourceNodeIndex, numberOfEpochs, learningRate, minimumDeliveryRatio, destinationPort, destinations,
					messageSize, numberOfPacketsPerLink, miniBatchSize, interval, secPathMabType, enableDadeAlgorithm,
					enableDedaAlgorithm, minAckForRttEstimation, experimentTimeElapsedMs)
				if err != nil {
					fmt.Printf("start osmd failed due to: %v", err)
				}
			}()
			return nil
		},
	}
	secPathMabBatchEvents = append(secPathMabBatchEvents, initOsmdEvent)

	// 设置规划好的破坏率改变序列
	currentTime += 10 * time.Second
	changeCorruptRatioEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_ChangeCorruptRatio,
		Handler: func() error {
			err = ChangeCorruptRatioInTimStampLevel()
			if err != nil {
				fmt.Printf("change corrupt ratio in timestamp level failed due to: %v", err)
			}
			return nil
		},
	}
	secPathMabBatchEvents = append(secPathMabBatchEvents, changeCorruptRatioEvent)

	// 3. 进行同步
	currentTime += 5 * time.Second
	syncEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_SynchronizeTimestamp,
		Handler: func() error {
			for i := range len(topology_manager.TopologyInstance.Nodes) {
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

	// 5. 进行结果的拷贝
	currentTime += time.Duration(experimentTimeElapsedMs/1000+10) * time.Second

	copyEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_ResultHandling,
		Handler: func() error {
			go func() {
				err = dir.CopyDir("/home/zhf/Projects/emulator/backend/simulation/LirNode-1/output",
					fmt.Sprintf("%s%d", "/home/zhf/Projects/emulator/backend/result/", currentExperimentIndex))
				if err != nil {
					fmt.Printf("copy dir failed due to: %v\n", err)
				}
			}()
			return nil
		},
	}
	secPathMabBatchEvents = append(secPathMabBatchEvents, copyEvent)

	// 6. 进行拓扑的清空
	currentTime += 20 * time.Second
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

func ChangeCorruptRatioInEpochLevel() error {
	var err error
	startEpoch := 100
	updateInterval := 100
	maxUpdateCount := 3
	currentUpdateCount := 0
	currentEpoch := startEpoch
	for {
		if currentUpdateCount%2 == 0 {
			err = validation_manager.SetScheduledMaliciousParams(5, currentEpoch,
				250000, 250000,
				0, 0)
			if err != nil {
				return fmt.Errorf("change corrupt ratio failed due to: %v", err)
			}
			err = validation_manager.SetScheduledMaliciousParams(6, currentEpoch,
				50000, 50000,
				0, 0)
			if err != nil {
				return fmt.Errorf("change corrupt ratio failed due to: %v", err)
			}
		} else {
			err = validation_manager.SetScheduledMaliciousParams(5, currentEpoch,
				50000, 50000,
				0, 0)
			if err != nil {
				return fmt.Errorf("change corrupt ratio failed due to: %v", err)
			}
			err = validation_manager.SetScheduledMaliciousParams(6, currentEpoch,
				250000, 250000,
				0, 0)
			if err != nil {
				return fmt.Errorf("change corrupt ratio failed due to: %v", err)
			}
		}

		currentUpdateCount += 1
		currentEpoch += updateInterval
		if currentUpdateCount >= maxUpdateCount {
			break
		}
	}
	return nil
}

func ChangeCorruptRatioInTimStampLevel() error {
	startTimestamp := 10000
	updateInterval := 600
	maxUpdateCount := 20
	currentUpdateCount := 0
	currentTimestamp := startTimestamp
	for {
		if currentUpdateCount%2 == 0 {
			err := validation_manager.SetScheduledMaliciousParams(5, currentTimestamp,
				300000, 300000,
				0, 0)
			if err != nil {
				return fmt.Errorf("change corrupt ratio failed due to: %v", err)
			}
			err = validation_manager.SetScheduledMaliciousParams(6, currentTimestamp,
				25000, 25000,
				0, 0)
			if err != nil {
				return fmt.Errorf("change corrupt ratio failed due to: %v", err)
			}
			// 295 (25000) 17:20:04 207468 (300000) 17:20:09 407697 (25000) 17:20:15
		} else {
			err := validation_manager.SetScheduledMaliciousParams(5, currentTimestamp,
				25000, 25000,
				0, 0)
			if err != nil {
				return fmt.Errorf("change corrupt ratio failed due to: %v", err)
			}
			err = validation_manager.SetScheduledMaliciousParams(6, currentTimestamp,
				300000, 300000,
				0, 0)
			if err != nil {
				return fmt.Errorf("change corrupt ratio failed due to: %v", err)
			}
		}

		currentUpdateCount += 1
		currentTimestamp += updateInterval
		if currentUpdateCount >= maxUpdateCount {
			break
		}
	}
	return nil
}

// SecPathMabExperiment 进行多次的实验
func SecPathMabExperiment() error {
	configurationSettings := []*entities.ConfigurationSetting{
		//{
		//	Index: 1,
		//	Mapping: map[string]string{
		//		"per_link_delay": "1",
		//	},
		//},
		// ------------------------------------------
		{
			Mapping: map[string]string{
				"per_link_delay":             "20",
				"min_ack_for_rtt_estimation": "100",
			},
		},
		//{
		//	Mapping: map[string]string{
		//		"per_link_delay":             "2.5",
		//		"number_of_packets_per_link": "150",
		//	},
		//},
		//{
		//	Mapping: map[string]string{
		//		"per_link_delay":             "2.5",
		//		"number_of_packets_per_link": "200",
		//	},
		//},
		//{
		//	Mapping: map[string]string{
		//		"per_link_delay":             "2.5",
		//		"number_of_packets_per_link": "250",
		//	},
		//},
		//{
		//	Mapping: map[string]string{
		//		"per_link_delay":             "2.5",
		//		"number_of_packets_per_link": "300",
		//	},
		//},
		//{
		//	Mapping: map[string]string{
		//		"per_link_delay":             "2.5",
		//		"number_of_packets_per_link": "400",
		//	},
		//},
		//{
		//	Mapping: map[string]string{
		//		"per_link_delay":             "2.5",
		//		"number_of_packets_per_link": "500",
		//	},
		//},
		//{
		//	Mapping: map[string]string{
		//		"per_link_delay":             "2.5",
		//		"number_of_packets_per_link": "600",
		//	},
		//},
		// ------------------------------------------
		//{
		//	Mapping: map[string]string{
		//		"per_link_delay":             "2.5",
		//		"min_ack_for_rtt_estimation": "100",
		//	},
		//},
		//{
		//	Mapping: map[string]string{
		//		"per_link_delay":             "2.5",
		//		"min_ack_for_rtt_estimation": "150",
		//	},
		//},
		//{
		//	Mapping: map[string]string{
		//		"per_link_delay":             "2.5",
		//		"min_ack_for_rtt_estimation": "200",
		//	},
		//},
		//{
		//	Mapping: map[string]string{
		//		"per_link_delay":             "2.5",
		//		"min_ack_for_rtt_estimation": "250",
		//	},
		//},
		//{
		//	Mapping: map[string]string{
		//		"per_link_delay":             "2.5",
		//		"min_ack_for_rtt_estimation": "300",
		//	},
		//},
		//{
		//	Mapping: map[string]string{
		//		"per_link_delay": "10",
		//	},
		//},
		//{
		//	Mapping: map[string]string{
		//		"per_link_delay": "20",
		//	},
		//},
	}
	for index, configurationSetting := range configurationSettings {
		experimentIndex := index + 1
		fmt.Printf("sec path mab experiment\n")
		secPathMabEvents, err := GenerateSecPathMabEvents(experimentIndex, configurationSetting)
		fmt.Printf("number of events: %d\n", len(secPathMabEvents))
		if err != nil {
			fmt.Printf("error to generate sec path mab batch events: %v\n", err)
		}

		err = experiments.SingleSimulation(configurationSetting, secPathMabEvents)
		if err != nil {
			return fmt.Errorf("sec path mab batch experiment failed: %w", err)
		}
	}

	return nil
}
