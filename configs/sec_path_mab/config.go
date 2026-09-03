package sec_path_mab

type Config struct {
	Enabled                   bool    `mapstructure:"enabled"`
	ExperimentType            int     `mapstructure:"experiment_type"`
	NumberOfHops              int     `mapstructure:"number_of_hops"`
	NumberOfIntermediateNodes int     `mapstructure:"number_of_intermediate_nodes"`
	LowCorruptRatio           float64 `mapstructure:"low_corrupt_ratio"`
	HighCorruptRatio          float64 `mapstructure:"high_corrupt_ratio"`
}
