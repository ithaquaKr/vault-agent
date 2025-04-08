package config

import (
	"github.com/ithaquaKr/vault-agent/internal/vault"
	"github.com/spf13/viper"
)

type Config struct {
	VaultConfig vault.VaultConfig `mapstructure:"vaultConfig"`
}

func LoadConfig(path, filename string) (*Config, error) {
	// Read configurations file
	v := viper.New()
	v.SetConfigFile(filename)
	v.AddConfigPath(path)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, err
		}
		return nil, err
	}
	// TODO: Templating data, support load configs from environment and add to YAML file.
	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
