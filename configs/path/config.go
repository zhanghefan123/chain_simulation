package path

type Config struct {
	Cmd              string `mapstructure:"cmd"`
	ConfigurationYml string `mapstructure:"configuration_yml"`
}
