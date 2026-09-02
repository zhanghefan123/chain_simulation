package entities

// StartServer processes: int = -1, listen_port: int = -1, server_type: str = "", interface: str = "", ip_version: str = ""
type StartServer struct {
	Processes            int    `json:"processes"`
	ListenPort           int    `json:"listen_port"`
	SimulationIndex      int    `json:"simulation_index"`
	NumberOfDestinations int    `json:"number_of_destinations"`
	SelectedNetworkLayer string `json:"selected_network_layer"`
	ServerType           string `json:"server_type"`
	Interface            string `json:"interface"`
	IpVersion            string `json:"ip_version"`
	ExperimentName       string `json:"experiment_name"`
}

func NewStartServer(processes int, selectedNetworkLayer string, listenPort int, simulationIndex int,
	serverType string, Interface string, IpVersion string, NumberOfDestinations int, ExperimentName string) *StartServer {
	return &StartServer{
		Processes:            processes,
		SelectedNetworkLayer: selectedNetworkLayer,
		ListenPort:           listenPort,
		SimulationIndex:      simulationIndex,
		ServerType:           serverType,
		Interface:            Interface,
		IpVersion:            IpVersion,
		NumberOfDestinations: NumberOfDestinations,
		ExperimentName:       ExperimentName,
	}
}
