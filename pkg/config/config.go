package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type VaultConfig struct {
	InitConfig VaultInitConfig `mapstructure:"init"`
	Data       VaultData       `mapstructure:"data"`
}

type VaultInitConfig struct {
	KeyShares int `mapstructure:"keyShares"`
	Threshold int `mapstructure:"threshold"`
}

type VaultData struct {
	Policies []Policy `mapstructure:"policies"`
}

type Policy struct {
	Name  string `mapstructure:"name"`
	Rules string `mapstructure:"rules"`
}

type Config struct {
	VaultConfig VaultConfig `mapstructure:"vault"`
}

func LoadConfig(v *viper.Viper) (*Config, error) {
	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

// BindFlag binds a cobra flag to a viper instance.
func BindFlag(cmd *cobra.Command, key string, defaultValue interface{}, description string, vp *viper.Viper) error {
	flags := cmd.PersistentFlags()

	switch v := defaultValue.(type) {
	case bool:
		flags.Bool(key, v, description)
	case int:
		flags.Int(key, v, description)
	case string:
		flags.String(key, v, description)
	case []string:
		flags.StringSlice(key, v, description)
	default:
		return fmt.Errorf("unsupported flag type for key: %s", key)
	}

	return vp.BindPFlag(key, flags.Lookup(key))
}
