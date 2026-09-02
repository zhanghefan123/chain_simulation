package validation_manager

import (
	"chain_simulation/configs"
	"chain_simulation/entities"
	"chain_simulation/entities/types"
	"chain_simulation/utils/request"
	"fmt"
	"time"
)

var (
	StartClientUrl                = "startClient"
	StartServerUrl                = "startServer"
	InitOsmdUrl                   = "initOsmd"
	StartOsmdUrl                  = "startOsmd"
	SetSchduledMaliciousParamsUrl = "setSchduledMaliciousParams"
	ModifyBloomFilterUrl          = "modifyBloomFilter"
	StartSyncUrl                  = "startSync"
	InsertSessionTableEntriesUrl  = "insertSessionTableEntries"
)

type ValidationManager struct{}

func StartClient(nodeIndex int, processes int, selectedNetworkLayer string, destinationPort int, destinations []string, transmissionPattern string, fileSize int, bufferSize int, content string,
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

func StartServer(nodeIndex int, processes int, selectedNetworkLayer string, simulationIndex int, listenPort int, serverType string, networkInterface string, ipVersion string, numberOfDestinations int, experimentName string) error {
	// 进行监听端口的获取
	containerListenPort := configs.TopConfigInstance.NetworkConfig.ValidationNodePort + nodeIndex
	// 进行 url 的构造
	startServerUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		containerListenPort,
		StartServerUrl)
	// 构造数据
	startServerInstance := entities.NewStartServer(processes, selectedNetworkLayer, listenPort, simulationIndex, serverType, networkInterface, ipVersion, numberOfDestinations, experimentName)
	// 进行 request
	err := request.PostJson(startServerUrl, startServerInstance)
	if err != nil {
		return fmt.Errorf("post start server failed %v", err)
	}
	return nil
}

func InitOsmd(nodeIndex int,
	learningRate float64, minimumDeliveryRatio float64,
	destinationPort int, destinations []string,
	messageSize int, numberOfPktsPerLink int, miniBatchSize int, packetSendingInterval float64, secPathMabStrategy types.SecPathMabStrategy,
	enableDadeAlgorithm bool, enableDedaAlgorithm bool, minAckForRttEstimation int, experimentTimeElapsedMs int) error {
	// 进行监听端口的获取
	containerListenPort := configs.TopConfigInstance.NetworkConfig.ValidationNodePort + nodeIndex
	// 进行 url 的构造
	initOsmdUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		containerListenPort,
		InitOsmdUrl)
	// 构造数据
	initOsmdInstance := entities.NewInitOsmd(learningRate, minimumDeliveryRatio,
		destinationPort, destinations, messageSize,
		numberOfPktsPerLink, miniBatchSize,
		packetSendingInterval, secPathMabStrategy, enableDadeAlgorithm, enableDedaAlgorithm,
		minAckForRttEstimation, experimentTimeElapsedMs)
	// 进行 request
	err := request.PostJson(initOsmdUrl, initOsmdInstance)
	if err != nil {
		return fmt.Errorf("post init osmd failed %v", err)
	}
	return nil
}

func StartOsmd(nodeIndex int) error {
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

func SetScheduledMaliciousParams(nodeIndex, employedEpochOrTimestampMs, corruptRatioStart, corruptRatioEnd, corruptSpecialPacketRatioStart, corruptSpecialPacketRatioEnd int) error {
	// 进行 container 的监听节点的获取
	containerListenPort := configs.TopConfigInstance.NetworkConfig.ValidationNodePort + nodeIndex
	// 进行 source 的监听节点的获取
	sourceListenPort := configs.TopConfigInstance.NetworkConfig.ValidationNodePort + 1
	// 进行 container url 的构造
	setContainerScheduledMaliciousParamsUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		containerListenPort,
		SetSchduledMaliciousParamsUrl)
	// 进行 source url 的构造
	setSourceScheduledMaliciousParamsUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		sourceListenPort,
		SetSchduledMaliciousParamsUrl)
	// 进行参数实例的构造
	scheduledMaliciousParmasInstance := entities.NewScheduledMaliciousParams(employedEpochOrTimestampMs, nodeIndex, corruptRatioStart,
		corruptRatioEnd, corruptSpecialPacketRatioStart, corruptSpecialPacketRatioEnd)
	// 利用两个 url 进行分别的 post
	err := request.PostJson(setContainerScheduledMaliciousParamsUrl, scheduledMaliciousParmasInstance)
	if err != nil {
		return fmt.Errorf("post container set scheduled malicious params failed %v", err)
	}
	err = request.PostJson(setSourceScheduledMaliciousParamsUrl, scheduledMaliciousParmasInstance)
	if err != nil {
		return fmt.Errorf("post source set scheduled malicious params failed %v", err)
	}
	return nil
}

func ModifyBloomFilter(nodeIndex int, bfEffectiveBits int) error {
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

func StartSync(nodeIndex int) error {
	// 进行监听端口的获取
	containerListenPort := configs.TopConfigInstance.NetworkConfig.ValidationNodePort + nodeIndex

	// 进行 url 的构造
	startSynchronizeUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		containerListenPort,
		StartSyncUrl)

	currentTimestamp := time.Now().UnixMilli()

	// 进行请求的发送
	err := request.PostJson(startSynchronizeUrl, &entities.SyncInstance{
		CurrentTimestamp: currentTimestamp,
	})
	if err != nil {
		return fmt.Errorf("post start synchronize url failed %v", err)
	}
	return nil
}

func InsertSessionTableEntries(nodeIndex int, numberOfEntries int) error {
	// 进行监听端口的获取
	containerListenPort := configs.TopConfigInstance.NetworkConfig.ValidationNodePort + nodeIndex

	// 进行 url 的构造
	insertSessionTableEntriesUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		containerListenPort,
		InsertSessionTableEntriesUrl)

	// 进行请求的发送
	err := request.PostJson(insertSessionTableEntriesUrl, &entities.InsertSessionTableEntriesInstance{
		NumberOfEntries: numberOfEntries,
	})
	if err != nil {
		return fmt.Errorf("post insert session table entries url failed %v", err)
	}
	return nil
}
