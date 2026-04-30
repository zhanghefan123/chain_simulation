package network

type Config struct {
	BackendAddr        string `mapstructure:"backend_addr"`
	BackendPort        int    `mapstructure:"backend_port"`
	ValidationNodePort int    `mapstructure:"validation_node_port"`
}
