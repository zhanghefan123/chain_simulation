package configs

import (
	"chain_simulation/configs/attack"
	"chain_simulation/configs/consensus"
	"chain_simulation/configs/network"
	"chain_simulation/configs/path"
	"chain_simulation/configs/url"
	"fmt"

	"github.com/spf13/viper"
)

var (
	ConfigurationFilePath = "../resources/configuration.yml"
	TopConfigInstance     = &TopConfig{}
)

type TopConfig struct {
	NetworkConfig   network.Config   `mapstructure:"network_config"`
	AttackConfig    attack.Config    `mapstructure:"attack_config"`
	PathConfig      path.Config      `mapstructure:"path_config"`
	UrlConfig       url.Config       `mapstructure:"url_config"`
	ConsensusConfig consensus.Config `mapstructure:"consensus_config"`
}

func InitTopConfig() error {
	// 定义错误
	var err error
	// 初始化解释器
	tempViper := viper.New()
	// 设置yml文件地址
	tempViper.SetConfigFile(ConfigurationFilePath)
	// 将yml加载
	if err = tempViper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
	}
	// 反序列化
	if err = tempViper.Unmarshal(TopConfigInstance); err != nil {
		return fmt.Errorf("unmarshal config error: %v", err)
	}
	return nil
}
