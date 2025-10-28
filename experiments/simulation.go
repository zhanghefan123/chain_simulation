package experiments

import (
	"chain_simulation/configs"
	"chain_simulation/entities"
	"chain_simulation/modules/backend_manager"
	"chain_simulation/modules/scheduler"
	"chain_simulation/modules/thread_manager"
	"chain_simulation/utils/file"
	"fmt"
	"time"
)

var experimentIndex = 0

// SingleSimulation 一轮实验
func SingleSimulation(configurationSetting *entities.ConfigurationSetting, events []*entities.Event) error {
	// 1. 进行修改 configuration.yml 的修改
	err := file.ModifyYml(configs.TopConfigInstance.PathConfig.ConfigurationYml, configurationSetting.Mapping)
	if err != nil {
		return fmt.Errorf("modify yml configuration failed: %v", err)
	}
	// 2. 进行后端的启动
	backend_manager.StartBackendService(experimentIndex)
	// 3. 根据拓扑选择执行的 event 序列
	scheduler.SetEventsIntoScheduler(events)
	// 4. 进行 scheduler 的启动
	scheduler.StartScheduler()
	// 5. 检查是否 simulation 已经结束
	ClearLastSimulation(configurationSetting.Mapping)
	// 6. 结束之后进行 wait
	thread_manager.ThreadManagerInstance.Wait()
	// 7. 更新 experiment index
	experimentIndex++
	return nil
}

// ClearLastSimulation 当一轮实验结束后需要将所有的环境进行清空
func ClearLastSimulation(mapping map[string]string) {
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
				fmt.Printf("simulation for mapping %v is finished \n", mapping)
				break ForLoop
			}
		}
	}
}
