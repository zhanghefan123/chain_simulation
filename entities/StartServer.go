package entities

// StartServer processes: int = -1, listen_port: int = -1, server_type: str = "", interface: str = "", ip_version: str = ""
type StartServer struct {
	Processes            int    `json:"processes"`
	SelectedNetworkLayer string `json:"selected_network_layer"`
	ListenPort           int    `json:"listen_port"`
	SimulationIndex      int    `json:"simulation_index"`
	ServerType           string `json:"server_type"`
	Interface            string `json:"interface"`
	IpVersion            string `json:"ip_version"`
	NumberOfDestinations int    `json:"number_of_destinations"`
}

func NewStartServer(processes int, selectedNetworkLayer string, listenPort int, simulationIndex int, serverType string, Interface string, IpVersion string, NumberOfDestinations int) *StartServer {
	return &StartServer{
		Processes:            processes,
		SelectedNetworkLayer: selectedNetworkLayer,
		ListenPort:           listenPort,
		SimulationIndex:      simulationIndex,
		ServerType:           serverType,
		Interface:            Interface,
		IpVersion:            IpVersion,
		NumberOfDestinations: NumberOfDestinations,
	}
}
