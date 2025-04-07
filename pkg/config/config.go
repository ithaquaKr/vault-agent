package config

import (
	"fmt"
	"log/slog"
	"os"
)

type ConfigFile struct {
	Path string
	Data map[string]interface{}
}

func parseConfiguration(vaultConfigFile string) *ConfigFile {
	// Read configurations file
	vaultConfig, err := os.ReadFile(vaultConfigFile)
	if err != nil {
		slog.Error(fmt.Sprintf("error reading config file template: %s", err))
	}

	// TODO: Templating data, support load configs from environment and add to YAML file.

	var data map[string]interface{}

	return &ConfigFile{
		Path: vaultConfigFile,
		Data: data,
	}
}
