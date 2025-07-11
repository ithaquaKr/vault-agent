package main

import (
	"fmt"
	"log/slog"

	"github.com/ithaquaKr/vault-agent/client"
	vaultManager "github.com/ithaquaKr/vault-agent/internal/vault"
	"github.com/ithaquaKr/vault-agent/pkg/config"
	"github.com/spf13/cobra"
)

const (
	cfgSecretShares    = "secret-shares"
	cfgSecretThreshold = "secret-threshold"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the target Vault instance",
	Long: `This command will verify the backend service is accessible, then
run "vault init" against the target Vault instance, before encrypting and
storing the keys in the given backend.

It will not unseal the Vault instance after initializing.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		slog.Info("Init Vault...")

		client, err := client.NewVaultClient(c.GetString(cfgVaultAddress), true)
		if err != nil {
			return fmt.Errorf("error connecting to Vault: %w", err)
		}

		v, err := vaultManager.New(client, cfg.VaultConfig.InitConfig, cfg.VaultConfig.Data)
		if err != nil {
			return fmt.Errorf("error creating Vault connect: %w", err)
		}

		if err = v.Init(); err != nil {
			return fmt.Errorf("error initializing Vault: %w", err)
		}

		return nil
	},
}

func init() {
	config.BindFlag(initCmd, cfgSecretShares, 5, "Number of keys that Vault will create during the initialization step.", c)
	config.BindFlag(initCmd, cfgSecretThreshold, 3, "Number of keys required to unseal Vault.", c)
}
