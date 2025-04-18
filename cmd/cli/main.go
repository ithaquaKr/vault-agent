package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Version = "dev"

const cfgVaultAddress = "vault-addr"

const cfgConfigPath = "config-path"

var c = viper.New()

var rootCmd = &cobra.Command{
	Use:     "vault-agent-cli",
	Short:   "Sync configurations, secrets of Hashicorp Vault",
	Version: Version, // TODO: Implement versioning via another command
	Long:    "This is a tool to help setup, management of Hashicorp Vault.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
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
		slog.Error(fmt.Sprintf("error executing command: %s", err.Error()))
		os.Exit(1)
	}
}

func init() {
	configVar(rootCmd, cfgVaultAddress, "http://127.0.0.1:8200", "The URL of the remote Vault", c)
	configVar(rootCmd, cfgConfigPath, "", "The configs file path", c)
}

func main() {
	execute()
}
