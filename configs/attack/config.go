package attack

type Config struct {
	Type        string `mapstructure:"type"`
	ThreadCount int    `mapstructure:"thread_count"`
	Duration    int    `mapstructure:"duration"`
}
