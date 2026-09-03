package serviceapi

import (
	"fmt"

	"chain_simulation/configs"
	"chain_simulation/utils/request"
)

// PostBackend sends a request to the main backend HTTP service.
func PostBackend(endpointPath string, payload any) error {
	return post(endpointURL(configs.TopConfigInstance.NetworkConfig.BackendPort, endpointPath), payload)
}

// PostBackendWithPortOffset sends a request to a backend-side service whose
// port is relative to the main backend port.
func PostBackendWithPortOffset(portOffset int, endpointPath string, payload any) error {
	port := configs.TopConfigInstance.NetworkConfig.BackendPort + portOffset
	return post(endpointURL(port, endpointPath), payload)
}

// PostValidationNode sends a request to one validation node. Validation node
// indexes are mapped to ports relative to ValidationNodePort.
func PostValidationNode(nodeIndex int, endpointPath string, payload any) error {
	port := configs.TopConfigInstance.NetworkConfig.ValidationNodePort + nodeIndex
	return post(endpointURL(port, endpointPath), payload)
}

func endpointURL(port int, endpointPath string) string {
	return fmt.Sprintf(
		"http://%s:%d/%s",
		configs.TopConfigInstance.NetworkConfig.BackendAddr,
		port,
		endpointPath,
	)
}

func post(url string, payload any) error {
	if err := request.PostJson(url, payload); err != nil {
		return fmt.Errorf("POST %s failed: %w", url, err)
	}
	return nil
}
