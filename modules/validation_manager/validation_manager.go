package validation_manager

import (
	"chain_simulation/configs"
	"chain_simulation/entities"
	"chain_simulation/entities/types"
	"chain_simulation/utils/request"
	"fmt"
)

var (
	ValidationManagerInstance     = &ValidationManager{}
	StartClientUrl                = "startClient"
	StartServerUrl                = "startServer"
	InitOsmdUrl                   = "initOsmd"
	StartOsmdUrl                  = "startOsmd"
	StartRetrieveAcksUrl          = "startRetrieveAcks"
	SetSchduledMaliciousParamsUrl = "setSchduledMaliciousParams"
	ModifyBloomFilterUrl          = "modifyBloomFilter"
	StartSyncUrl                  = "startSync"
)

type ValidationManager struct{}

func StartClient(nodeIndex int, processes int, selectedNetworkLayer string, destinationPort int, destinations []string, transmissionPattern string, fileSize int, bufferSize int, content string,
	batchSize int, messageSize int, interval float64) error {
	err := ValidationManagerInstance.StartClientInner(nodeIndex, processes, selectedNetworkLayer, destinationPort, destinations,
		transmissionPattern, fileSize, bufferSize, content,
		batchSize, messageSize, interval)
	if err != nil {
		return fmt.Errorf("start client inner failed: %v", err)
	} else {
		return nil
	}
}

func StartServer(nodeIndex int, processes int, selectedNetworkLayer string, simulationIndex int, listenPort int, serverType string, networkInterface string, ipVersion string, numberOfDestinations int) error {
	err := ValidationManagerInstance.StartServerInner(nodeIndex, processes, selectedNetworkLayer, simulationIndex, listenPort, serverType, networkInterface, ipVersion, numberOfDestinations)
	if err != nil {
		return fmt.Errorf("start server inner failed: %v", err)
	} else {
		return nil
	}
}

func InitOsmd(nodeIndex int, numberOfEpochs int,
	learningRate float64, minimumDeliveryRatio float64,
	destinationPort int, destinations []string,
	messageSize int, numberOfPktsPerLink int, miniBatchSize int, packetSendingInterval float64, secPathMabStrategy types.SecPathMabStrategy,
	enableDadeAlgorithm bool, enableDedaAlgorithm bool, minAckForRttEstimation int, experimentTimeElapsedMs int) error {
	err := ValidationManagerInstance.InitOsmdInner(nodeIndex, numberOfEpochs, learningRate, minimumDeliveryRatio,
		destinationPort, destinations, messageSize, numberOfPktsPerLink, miniBatchSize, packetSendingInterval, secPathMabStrategy,
		enableDadeAlgorithm, enableDedaAlgorithm, minAckForRttEstimation, experimentTimeElapsedMs)
	if err != nil {
		return fmt.Errorf("init osmd failed: %v", err)
	} else {
		return nil
	}
}

func StartOsmd(nodeIndex int) error {
	err := ValidationManagerInstance.StartOsmdInner(nodeIndex)
	if err != nil {
		return fmt.Errorf("start osmd failed: %v", err)
	} else {
		return nil
	}
}

func SetScheduledMaliciousParams(nodeIndex, employedEpochOrTimestampMs, corruptRatioStart, corrutpRatioEnd, corruptSpecialPacketRatioStart, corruptSpecialPacketRatioEnd int) error {
	err := ValidationManagerInstance.SetScheduledMaliciousParamsInner(nodeIndex, employedEpochOrTimestampMs, corruptRatioStart, corrutpRatioEnd, corruptSpecialPacketRatioStart, corruptSpecialPacketRatioEnd)
	if err != nil {
		return fmt.Errorf("set scheduled malicious params inner failed: %v", err)
	}
	err = ValidationManagerInstance.SetScheduledMaliciousParamsToSource(nodeIndex, employedEpochOrTimestampMs, corruptRatioStart, corrutpRatioEnd, corruptSpecialPacketRatioStart, corruptSpecialPacketRatioEnd)
	if err != nil {
		return fmt.Errorf("set scheduled malicious params source failed: %v", err)
	}
	return nil
}

func StartRetrieveAcks(nodeIndex int) error {
	err := ValidationManagerInstance.StartRetrieveAcksInner(nodeIndex)
	if err != nil {
		return fmt.Errorf("start retrieve acks inner failed: %v", err)
	} else {
		return nil
	}
}

func ModifyBloomFilter(nodeIndex int, bfEffectiveBits int) error {
	err := ValidationManagerInstance.ModifyBloomFilter(nodeIndex, bfEffectiveBits)
	if err != nil {
		return fmt.Errorf("modify bloom filter failed: %v", err)
	} else {
		return nil
	}
}

func (vm *ValidationManager) ModifyBloomFilter(nodeIndex int, bfEffectiveBits int) error {
	// 进行节点监听端口的获取
	listenPort := configs.TopConfigInstance.NetworkConfig.ValidationNodePort + nodeIndex
	// 进行 url 构造
	modifyBloomFilterUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		listenPort,
		ModifyBloomFilterUrl)
	// 构造数据
	modifyBloomFilterInstance := entities.NewModifyBloomFilter(bfEffectiveBits)
	// 进行 request
	err := request.PostJson(modifyBloomFilterUrl, modifyBloomFilterInstance)
	if err != nil {
		return fmt.Errorf("post modify bloom filter failed: %v", err)
	}
	return nil
}

func (vm *ValidationManager) StartClientInner(nodeIndex int, processes int, selectedNetworkLayer string, destinationPort int, destinations []string,
	transmissionPattern string, fileSize int, bufferSize int, content string,
	batchSize int, messageSize int, interval float64) error {
	// 进行节点的监听端口的获取
	listenPort := configs.TopConfigInstance.NetworkConfig.ValidationNodePort + nodeIndex
	// 进行 url 的构造
	startClientUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		listenPort,
		StartClientUrl)
	// 构造数据
	startClientInstance := entities.NewStartClient(selectedNetworkLayer, processes, destinationPort, destinations, transmissionPattern,
		fileSize, bufferSize, content, batchSize, messageSize, interval)
	// 进行 request
	err := request.PostJson(startClientUrl, startClientInstance)
	if err != nil {
		return fmt.Errorf("post start client failed %v", err)
	}
	return nil
}

func (vm *ValidationManager) StartServerInner(nodeIndex int, processes int, selectedNetworkLayer string, simulationIndex int, listenPort int, serverType string, networkInterface string, ipVersion string, numberOfDestinations int) error {
	// 进行监听端口的获取
	containerListenPort := configs.TopConfigInstance.NetworkConfig.ValidationNodePort + nodeIndex
	// 进行 url 的构造
	startServerUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		containerListenPort,
		StartServerUrl)
	// 构造数据
	startServerInstance := entities.NewStartServer(processes, selectedNetworkLayer, listenPort, simulationIndex, serverType, networkInterface, ipVersion, numberOfDestinations)
	// 进行 request
	err := request.PostJson(startServerUrl, startServerInstance)
	if err != nil {
		return fmt.Errorf("post start server failed %v", err)
	}
	return nil
}

func (vm *ValidationManager) InitOsmdInner(nodeIndex int, numberOfEpochs int,
	learningRate float64, minimumDeliveryRatio float64,
	destinationPort int, destinations []string,
	messageSize int, numberOfPktsPerLink int, miniBatchSize int, packetSendingInterval float64,
	strategy types.SecPathMabStrategy, enableDadeAlgorithm bool, enableDedaAlgorithm bool, minAckForRttEstimation int, experimentTimeElapsedMs int) error {
	// 进行监听端口的获取
	containerListenPort := configs.TopConfigInstance.NetworkConfig.ValidationNodePort + nodeIndex
	// 进行 url 的构造
	initOsmdUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		containerListenPort,
		InitOsmdUrl)
	// 构造数据
	initOsmdInstance := entities.NewInitOsmd(numberOfEpochs, learningRate, minimumDeliveryRatio,
		destinationPort, destinations, messageSize,
		numberOfPktsPerLink, miniBatchSize,
		packetSendingInterval, strategy, enableDadeAlgorithm, enableDedaAlgorithm, minAckForRttEstimation, experimentTimeElapsedMs)
	// 进行 request
	err := request.PostJson(initOsmdUrl, initOsmdInstance)
	if err != nil {
		return fmt.Errorf("post init osmd failed %v", err)
	}
	return nil
}

func (vm *ValidationManager) StartOsmdInner(nodeIndex int) error {
	// 进行监听端口的获取
	containerListenPort := configs.TopConfigInstance.NetworkConfig.ValidationNodePort + nodeIndex
	// 进行 url 的构造
	startOsmdUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		containerListenPort,
		StartOsmdUrl)
	// 进行 request
	err := request.PostJson(startOsmdUrl, nil)
	if err != nil {
		return fmt.Errorf("pos start osmd failed %v", err)
	}
	return nil
}

func (vm *ValidationManager) StartRetrieveAcksInner(nodeIndex int) error {
	// 进行监听接口获取
	containerListenPort := configs.TopConfigInstance.NetworkConfig.ValidationNodePort + nodeIndex
	// 进行 url 的构造
	startRetrieveAcksUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		containerListenPort,
		StartRetrieveAcksUrl)
	// 进行 request
	err := request.PostJson(startRetrieveAcksUrl, nil)
	if err != nil {
		return fmt.Errorf("post start retrieve acks failed %v", err)
	}
	return nil
}

func (vm *ValidationManager) SetScheduledMaliciousParamsInner(nodeIndex, employedEpochOrTimestampMs, corruptRatioStart, corruptRatioEnd, corruptSpecialPacketRatioStart, corruptSpecialPacketRatioEnd int) error {
	// 进行监听接口的获取
	containerListenPort := configs.TopConfigInstance.NetworkConfig.ValidationNodePort + nodeIndex
	// 进行 url 的构造
	setScheduledMaliciousParamsUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		containerListenPort,
		SetSchduledMaliciousParamsUrl)
	// 进行 changeCorruptRatioInstance 的构造
	scheduledMaliciousParmasInstance := entities.NewScheduledMaliciousParams(employedEpochOrTimestampMs, nodeIndex, corruptRatioStart,
		corruptRatioEnd, corruptSpecialPacketRatioStart, corruptSpecialPacketRatioEnd)
	// 进行 request
	err := request.PostJson(setScheduledMaliciousParamsUrl, scheduledMaliciousParmasInstance)
	if err != nil {
		return fmt.Errorf("post set scheduled malicious params failed %v", err)
	}
	return nil
}

func (vm *ValidationManager) SetScheduledMaliciousParamsToSource(nodeIndex, employedEpochOrTimestampMs, corruptRatioStart, corruptRatioEnd, corruptSpecialPacketRatioStart, corruptSpecialPacketRatioEnd int) error {
	// 进行监听接口的获取
	sourceListenPort := configs.TopConfigInstance.NetworkConfig.ValidationNodePort + 1
	// 进行 url 的构造
	setScheduledMaliciousParamsUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		sourceListenPort,
		SetSchduledMaliciousParamsUrl)
	// 进行 changeCorruptRatioInstance
	scheduledMaliciousParmasInstance := entities.NewScheduledMaliciousParams(employedEpochOrTimestampMs, nodeIndex, corruptRatioStart,
		corruptRatioEnd, corruptSpecialPacketRatioStart, corruptSpecialPacketRatioEnd)
	// 进行 request
	err := request.PostJson(setScheduledMaliciousParamsUrl, scheduledMaliciousParmasInstance)
	if err != nil {
		return fmt.Errorf("post set scheduled malicious params failed %v", err)
	}
	return nil
}

func StartSync(nodeIndex int) error {
	err := ValidationManagerInstance.StartSyncInner(nodeIndex)
	if err != nil {
		return fmt.Errorf("start synchronize failed: %v", err)
	} else {
		return nil
	}
}

func (vm *ValidationManager) StartSyncInner(nodeIndex int) error {
	// 进行监听端口的获取
	containerListenPort := configs.TopConfigInstance.NetworkConfig.ValidationNodePort + nodeIndex

	// 进行 url 的构造
	startSynchronizeUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		containerListenPort,
		StartSyncUrl)

	err := request.PostJson(startSynchronizeUrl, nil)
	if err != nil {
		return fmt.Errorf("post start synchronize url failed %v", err)
	}
	return nil
}
