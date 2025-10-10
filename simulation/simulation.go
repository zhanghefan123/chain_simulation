package simulation

import (
	"chain_simulation/configs"
	"chain_simulation/entities/types"
	"chain_simulation/modules/backend_manager"
	"chain_simulation/modules/scheduler"
	"chain_simulation/modules/thread_manager"
	"chain_simulation/utils/file"
	"fmt"
	"time"
)

// SingleSimulation 一轮实验
func SingleSimulation(key, value string, topologyType types.TopologyType, experimentIndex int) error {
	// 1. 进行修改 configuration.yml 的修改
	err := file.ModifyYml(configs.TopConfigInstance.PathConfig.ConfigurationYml, key, value)
	if err != nil {
		return fmt.Errorf("modify yml configuration failed: %v", err)
	}
	// 2. 进行后端的启动
	backend_manager.StartBackendService(experimentIndex)
	// 3. 根据拓扑选择执行的 event 序列
	switch topologyType {
	case types.TopologyType_HyperledgerFabric:
		scheduler.SetEventsIntoScheduler(FabricEvents)
	case types.TopologyType_FiscoBcos:
		scheduler.SetEventsIntoScheduler(FiscoBcosEvents)
	case types.TopologyType_ChainMaker:
		scheduler.SetEventsIntoScheduler(ChainMakerEvents)
	}
	// 4. 进行 scheduler 的启动
	scheduler.StartScheduler()
	// 5. 检查是否 simulation 已经结束
	ClearLastSimulation(key, value)
	// 6. 结束之后进行 wait
	thread_manager.ThreadManagerInstance.Wait()

	return nil
}

// ClearLastSimulation 当一轮实验结束后需要将所有的环境进行清空
func ClearLastSimulation(key, value string) {
	// 1. 检查是否已经没有 event list 了
	var ticker *time.Ticker
	ticker = time.NewTicker(time.Second)
	defer ticker.Stop()
ForLoop:
	for {
		select {
		case <-ticker.C:
			if len(scheduler.SchedulerInstance.EventList) == 0 {
				backend_manager.StopBackendService()
				scheduler.StopScheduler()
				fmt.Printf("simulation for key: %v, value: %v finished \n", key, value)
				break ForLoop
			}
		}
	}
}
