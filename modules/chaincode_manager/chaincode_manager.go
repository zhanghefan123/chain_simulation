package chaincode_manager

import (
	"chain_simulation/configs"
	"chain_simulation/utils/request"
	"fmt"
)

var ChainCodeManagerInstance = &ChainCodeManager{}

type ChainCodeManager struct {
}

func InstallChainCode() error {
	installChainCodeUrl := fmt.Sprintf("http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		configs.TopConfigInstance.NetworkConfig.BackendPort,
		configs.TopConfigInstance.UrlConfig.InstallChainCodeUrl)
	err := request.PostJson(installChainCodeUrl, nil)
	if err != nil {
		return fmt.Errorf("install chaincode error: %s", err)
	}
	return nil
}
