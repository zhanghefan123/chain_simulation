package attack_manager

import (
	"fmt"

	"chain_simulation/configs"
	"chain_simulation/entities"
	"chain_simulation/entities/types"
	"chain_simulation/modules/internal/serviceapi"
)

// StartAttack 对实验暴露的攻击接口
func StartAttack(topologyType types.TopologyType, attackDuration int, attackThreadCount int) error {
	attackInstance := entities.NewAttackInstance(topologyType, attackDuration, attackThreadCount)
	if err := serviceapi.PostBackendWithPortOffset(
		1,
		configs.TopConfigInstance.UrlConfig.StartAttackUrl,
		attackInstance,
	); err != nil {
		return fmt.Errorf("start attack: %w", err)
	}
	return nil
}
