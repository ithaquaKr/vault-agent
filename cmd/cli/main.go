package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "vault-agent-cli",
	Short: "Sync configurations, secrets of Hashicorp Vault",
	Long:  "This is a tool to help setup, management of Hashicorp Vault.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		cmd.SetContext(ctx)
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
	// TODO: Parse configurations
}

func main() {
	execute()
}
