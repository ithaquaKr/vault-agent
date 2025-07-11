package main

import (
	"fmt"
	"log/slog"

	"github.com/ithaquaKr/vault-agent/client"
	vaultManager "github.com/ithaquaKr/vault-agent/internal/vault"
	"github.com/spf13/cobra"
)

var unsealCmd = &cobra.Command{
	Use:   "unseal",
	Short: "Unseal the target Vault instance",
	Long: `Unseals Vault with unseal keys provide from command line.

  It will continuously attempt to unseal the target Vault instance, 
  by retrieving unseal keys from command line.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		slog.Info("Unsealing Vault...")

		client, err := client.NewVaultClient(c.GetString(cfgVaultAddress), true)
		if err != nil {
			return fmt.Errorf("error connecting to Vault: %w", err)
		}

		v, err := vaultManager.New(client, cfg.VaultConfig.InitConfig, cfg.VaultConfig.Data)
		if err != nil {
			return fmt.Errorf("error creating Vault manager: %w", err)
		}

		if err = v.Unseal(); err != nil {
			return fmt.Errorf("error unsealing Vault: %w", err)
		}

		return nil
	},
}
