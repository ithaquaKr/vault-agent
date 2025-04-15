package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ithaquaKr/vault-agent/client"
	vaultController "github.com/ithaquaKr/vault-agent/internal/vault"
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
	Run: func(_ *cobra.Command, _ []string) {
		slog.Info("Init Vault...")
		client, err := client.NewVaultClient("http://127.0.0.1:8200", true)
		if err != nil {
			slog.Error(fmt.Sprintf("error connecting to Vault: %s", err.Error()))
		}
		config, err := config.LoadConfig("../../", "config.yaml")
		if err != nil {
			slog.Error(fmt.Sprintf("can load configuration: %s", err))
		}
		fmt.Println(config.VaultConfig)

		v, err := vaultController.New(client, config.VaultConfig)
		if err != nil {
			slog.Error(fmt.Sprintf("error creating Vault connect: %s", err.Error()))
		}

		if err = v.Init(); err != nil {
			slog.Error(fmt.Sprintf("error initializing Vault: %s", err.Error()))
			os.Exit(1)
		}
	},
}

func init() {
	configVar(initCmd, cfgSecretShares, 5, "Number of keys that Vault will create during the initialization step.", c)
	configVar(initCmd, cfgSecretThreshold, 3, "Number of keys required to unseal Vault.", c)

	rootCmd.AddCommand(initCmd)
}
