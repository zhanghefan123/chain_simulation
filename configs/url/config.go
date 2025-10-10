package url

type Config struct {
	StartTopologyUrl    string `mapstructure:"start_topology_url"`
	StopTopologyUrl     string `mapstructure:"stop_topology_url"`
	StartAttackUrl      string `mapstructure:"start_attack_url"`
	StartTxRateTestUrl  string `mapstructure:"start_tx_rate_test_url"`
	StopTxRateTestUrl   string `mapstructure:"stop_tx_rate_test_url"`
	InstallChainCodeUrl string `mapstructure:"install_chain_code_url"`
}
