package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ithaquaKr/vault-agent/client"
	vaultServer "github.com/ithaquaKr/vault-agent/internal/vault"
	"github.com/ithaquaKr/vault-agent/pkg/config"
	"github.com/spf13/cobra"
)

var unsealCmd = &cobra.Command{
	Use:   "unseal",
	Short: "Unseal the target Vault instance",
	Long: `Unseals Vault with unseal keys provide from command line.
  It will continuously attempt to unseal the target Vault instance, by retrieving
  unseal keys from command line.`,
	Run: func(_ *cobra.Command, _ []string) {
		slog.Info("Unsealing Vault ...")
		// TODO: Move logic init VaultClient to some common function
		config, err := config.LoadConfig("../../", "config.yaml")
		if err != nil {
			slog.Error(fmt.Sprintf("can load configuration: %s", err))
		}
		fmt.Println(config.VaultConfig)

		client, err := client.NewVaultClient("http://127.0.0.1:8200", true)
		if err != nil {
			slog.Error(fmt.Sprintf("error connecting to Vault: %s", err.Error()))
		}
		v, err := vaultServer.New(client, config.VaultConfig)
		if err != nil {
			slog.Error(fmt.Sprintf("error creating Vault connect: %s", err.Error()))
		}

		keys := []string{"key1", "key2", "key3"}
		addresses := []string{"addresses1", "addresses2", "addresses3"}

		if err = v.Unseal(keys, addresses); err != nil {
			slog.Error(fmt.Sprintf("error unsealing Vault: %s", err))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(unsealCmd)
}
