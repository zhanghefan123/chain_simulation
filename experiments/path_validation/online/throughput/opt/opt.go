package opt

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

// 依然使用我们的 secPathMab 生成的 topology
var topologyType = types.TopologyType_SecPathMab

func GenerateOptEvents(setting *entities.ConfigurationSetting) ([]*entities.Event, error) {
	optEvents := make([]*entities.Event, 0)

	// 首先进行环境名称的获取
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
	optEvents = append(optEvents, installKernelEvent)

	// 进行拓扑的启动
	currentTime += time.Second * 10
	perLinkDelayInMs, err := strconv.ParseFloat(setting.Mapping["per_link_delay"], 64)
	if err != nil {
		return nil, fmt.Errorf("parse per link delay in ms error: %v", err)
	}
	startTopologyEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_StartTopology,
		Handler: func() error {
			return topology_manager.StartTopology(topologyType,
				&entities.DynamicParameters{
					PerLinkDelay: perLinkDelayInMs,
				})
		},
	}
	optEvents = append(optEvents, startTopologyEvent)

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

	// 添加服务器事件
	currentTime += time.Second * 10
	currentProcess := 1
	pathValidationProtocol := "OPT"
	simulationIndex := 1
	networkInterface := fmt.Sprintf("ln%d_idx1", serverIndex)
	serverEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_StartServer,
		Handler: func() error {
			go func() {
				err = validation_manager.StartServer(serverIndex, currentProcess, pathValidationProtocol, simulationIndex,
					31313, "multiprocess_file", networkInterface,
					"Ipv4", 1, finalExperimentName)
				if err != nil {
					fmt.Printf("epic start server error: %v", err)
				}
			}()
			return nil
		},
	}
	optEvents = append(optEvents, serverEvent)

	var destinations []string
	if destinationStr, ok := setting.Mapping["server_name"]; ok {
		destinations = append(destinations, destinationStr)
	}

	// 添加客户端事件
	currentTime += time.Second * 5
	clientEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_StartClient,
		Handler: func() error {
			go func() {
				err = validation_manager.StartClient(1, currentProcess, pathValidationProtocol, 31313, destinations,
					"file", 256, 512, "",
					0, 0, 0)
				if err != nil {
					fmt.Printf("start client error: %v", err)
				}
			}()
			return nil
		},
	}
	optEvents = append(optEvents, clientEvent)

	// 过一段时间之后进行拓扑的清空
	// 6. 进行拓扑的清空
	currentTime += 60 * time.Second
	stopEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_StopTopology,
		Handler: func() error {
			return topology_manager.StopTopology()
		},
	}
	optEvents = append(optEvents, stopEvent) // 每个节点发送 50 个 , 50 * 3 = 150, 最后一个节点 rtt = 12ms, 0.012 * 25000pkts/s = 300,  大约就是 450 个packet

	// 7. 进行内核模块的卸载
	currentTime += 10 * time.Second
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
	optEvents = append(optEvents, removeKernelEvent)

	// 8. 等待整个拓扑被清空
	currentTime += 20 * time.Second
	waitStopEvent := &entities.Event{
		StartTime: currentTime,
		Action:    types.ActionType_WaitTopologyRemove,
		Handler: func() error {
			return nil
		},
	}
	optEvents = append(optEvents, waitStopEvent)
	return optEvents, nil
}

func GenerateConfigurationSettings() ([]*entities.ConfigurationSetting, error) {
	// 最终的结果
	var configurationSettings []*entities.ConfigurationSetting
	// 选择每次 run 多少次
	runExperimentTimes := 10
	// 首先查询自己已经做过的实验
	alreadyRunedExperiments, err := breakpoint_awareness.GetAlreadyCaclulatedFileTransmissionResult("/home/zhf/Projects/emulator/backend/cmd/final_result/")
	if err != nil {
		return nil, fmt.Errorf("get already calculated file transmission result failed: %v", err)
	}
	// 拓扑参数设置 number_of_hops = 2,4,6,8,10, number_of_segments_per_hop = 1
	numberOfSegmentsPerHop := 1 // 注意在 fixed batch mode 之中, 这个只能为1
	for currentHop := 2; currentHop <= 10; currentHop += 2 {
		for runExperimentIndex := range runExperimentTimes {
			serverIndex := currentHop*(numberOfSegmentsPerHop+1) + 1
			experimentDestination := fmt.Sprintf("LirNode-%d", serverIndex)
			experimentName := fmt.Sprintf("PROTO:OPT|HOP:%d|SEG:%d||index:%d", currentHop,
				numberOfSegmentsPerHop, runExperimentIndex)
			if _, ok := alreadyRunedExperiments[experimentName]; ok {
				continue
			} else {
				configurationSetting := &entities.ConfigurationSetting{
					Mapping: map[string]string{
						"per_link_delay":             "0.0",
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
	return configurationSettings, nil
}

func RunExperiments() error {
	configurationSettings, err := GenerateConfigurationSettings()
	if err != nil {
		return fmt.Errorf("generate configuration settings failed: %v", err)
	}

	for _, configurationSetting := range configurationSettings {
		var optEvents []*entities.Event
		optEvents, err = GenerateOptEvents(configurationSetting)
		err = experiments.SingleSimulation(configurationSetting, optEvents)
		if err != nil {
			return fmt.Errorf("opt experiment failed: %v", err)
		}
	}

	return nil
}
