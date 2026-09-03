package validation_manager

import "chain_simulation/entities"

// StartClient starts a traffic client on the selected validation node.
func StartClient(nodeIndex int, params *entities.StartClient) error {
	return postToNode(nodeIndex, startClientPath, params, "start client")
}

// StartServer starts a traffic server on the selected validation node.
func StartServer(nodeIndex int, params *entities.StartServer) error {
	return postToNode(nodeIndex, startServerPath, params, "start server")
}
