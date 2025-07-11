package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ithaquaKr/vault-agent/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Version = "dev"

const cfgVaultAddress = "vault-addr"

const cfgConfigPath = "config-path"

var c = viper.New()

var cfg *config.Config

var rootCmd = &cobra.Command{
	Use:     "vault-agent-cli",
	Short:   "Sync configurations, secrets of Hashicorp Vault",
	Version: Version,
	Long:    "This is a tool to help setup, management of Hashicorp Vault.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		configPath := c.GetString(cfgConfigPath)
		if configPath != "" {
			c.SetConfigFile(configPath)
			if err := c.ReadInConfig(); err != nil {
				return fmt.Errorf("failed to read config file: %w", err)
			}
		}

		var err error
		cfg, err = config.LoadConfig(c)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		return nil
	},
}

func execute() {
	// Handle signal to prevent bad exit codes on "docker stop"
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT)
	go func() {
		<-sigs
		os.Exit(0)
	}()

	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func init() {
	config.BindFlag(rootCmd, cfgVaultAddress, "http://127.0.0.1:8200", "The URL of the remote Vault", c)
	config.BindFlag(rootCmd, cfgConfigPath, "", "The configs file path", c)

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(unsealCmd)
	rootCmd.AddCommand(syncCmd)
}

func main() {
	execute()
}
