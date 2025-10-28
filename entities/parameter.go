package entities

type ConfigurationSetting struct {
	Mapping map[string]string
}

func NewConfigurationSetting() *ConfigurationSetting {
	return &ConfigurationSetting{
		make(map[string]string),
	}
}
