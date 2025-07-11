package main

import (
	"fmt"
	"log/slog"

	"github.com/ithaquaKr/vault-agent/client"
	vaultManager "github.com/ithaquaKr/vault-agent/internal/vault"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize Vault configurations",
	Long:  `This command synchronizes policies and other configurations with the target Vault instance.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		slog.Info("Synchronizing Vault configurations...")

		client, err := client.NewVaultClient(c.GetString(cfgVaultAddress), true)
		if err != nil {
			return fmt.Errorf("error connecting to Vault: %w", err)
		}

		v, err := vaultManager.New(client, cfg.VaultConfig.InitConfig, cfg.VaultConfig.Data)
		if err != nil {
			return fmt.Errorf("error creating Vault manager: %w", err)
		}

		if err = v.SyncPolicy(); err != nil {
			return fmt.Errorf("error synchronizing policies: %w", err)
		}

		slog.Info("Vault synchronization completed successfully.")
		return nil
	},
}

