package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "vault-agent-cli",
	Short: "Sync configurations, secrets of Hashicorp Vault",
	Long:  "This is a tool to help setup, management of Hashicorp Vault.",
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error(fmt.Sprintf("error executing command: %s", err.Error()))
		os.Exit(1)
	}
}

func main() {
	execute()
}
