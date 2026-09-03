package chaincode_manager

import (
	"fmt"

	"chain_simulation/configs"
	"chain_simulation/modules/internal/serviceapi"
)

func InstallChainCode() error {
	if err := serviceapi.PostBackend(configs.TopConfigInstance.UrlConfig.InstallChainCodeUrl, nil); err != nil {
		return fmt.Errorf("install chaincode: %w", err)
	}
	return nil
}
