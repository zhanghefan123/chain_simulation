package validation_manager

import (
	"fmt"
	"time"

	"chain_simulation/entities"
)

// InitOsmd configures OSMD on the selected validation node.
func InitOsmd(nodeIndex int, params *entities.InitOsmd) error {
	return postToNode(nodeIndex, initOsmdPath, params, "initialize OSMD")
}

// StartOsmd starts a previously configured OSMD instance.
func StartOsmd(nodeIndex int) error {
	return postToNode(nodeIndex, startOsmdPath, nil, "start OSMD")
}

// SetScheduledMaliciousParams sends the same schedule to the affected node and
// to the source node, which coordinates the experiment.
func SetScheduledMaliciousParams(params *entities.ScheduledMaliciousParams) error {
	if params == nil {
		return fmt.Errorf("set scheduled malicious parameters: params must not be nil")
	}
	if err := postToNode(
		params.NodeId,
		setScheduledMaliciousParamsPath,
		params,
		"set scheduled malicious parameters on target node",
	); err != nil {
		return err
	}
	return postToNode(
		1,
		setScheduledMaliciousParamsPath,
		params,
		"set scheduled malicious parameters on source node",
	)
}

// StartSync synchronizes one validation node using the current local time.
func StartSync(nodeIndex int) error {
	params := &entities.SyncInstance{CurrentTimestamp: time.Now().UnixMilli()}
	return postToNode(nodeIndex, startSyncPath, params, "synchronize node time")
}
