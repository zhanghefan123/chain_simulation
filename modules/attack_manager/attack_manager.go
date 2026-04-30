package attack_manager

import (
	"chain_simulation/configs"
	"chain_simulation/entities"
	"chain_simulation/entities/types"
	"chain_simulation/utils/request"
	"fmt"
)

var AttackManagerInstance = &AttackManager{}

type AttackManager struct {
}

// StartAttack 对实验暴露的攻击接口
func StartAttack(topologyType types.TopologyType, attackDuration int, attackTheadCount int) error {
	err := AttackManagerInstance.Start(topologyType, attackDuration, attackTheadCount)
	if err != nil {
		return fmt.Errorf("start attack failed")
	}
	return nil
}

func (am *AttackManager) Start(topologyType types.TopologyType, attackDuration int, attackThreadCount int) error {
	// 因为这里只有一个恶意的节点
	maliciousPort := configs.TopConfigInstance.NetworkConfig.BackendPort + 1
	// 构造 url
	startAttackUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		maliciousPort,
		configs.TopConfigInstance.UrlConfig.StartAttackUrl)
	// 构造 Attack data
	attackInstance := entities.NewAttackInstance(topologyType, attackDuration, attackThreadCount)
	// 进行 request
	err := request.PostJson(startAttackUrl, attackInstance)
	if err != nil {
		return fmt.Errorf("post attack failed: %v", err)
	}
	return nil
}
