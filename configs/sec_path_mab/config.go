package sec_path_mab

type Config struct {
	NumberOfHops     int     `mapstructure:"number_of_hops"`
	LowCorruptRatio  float64 `mapstructure:"low_corrupt_ratio"`
	HighCorruptRatio float64 `mapstructure:"high_corrupt_ratio"`
}
