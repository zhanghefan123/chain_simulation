package validation_manager

import (
	"fmt"

	"chain_simulation/modules/internal/serviceapi"
)

const (
	startClientPath                 = "startClient"
	startServerPath                 = "startServer"
	initOsmdPath                    = "initOsmd"
	startOsmdPath                   = "startOsmd"
	setScheduledMaliciousParamsPath = "setSchduledMaliciousParams"
	modifyBloomFilterPath           = "modifyBloomFilter"
	startSyncPath                   = "startSync"
	insertSessionTableEntriesPath   = "insertSessionTableEntries"
)

func postToNode(nodeIndex int, endpointPath string, payload any, operation string) error {
	if err := serviceapi.PostValidationNode(nodeIndex, endpointPath, payload); err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return nil
}
