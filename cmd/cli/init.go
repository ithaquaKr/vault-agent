package main

import (
	"fmt"
	"log/slog"
	"os"

	vaultClient "github.com/hashicorp/vault-client-go"
	customVault "github.com/ithaquaKr/vault-agent/internal/vault"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the target Vault instance",
	Long: `This command will verify the backend service is accessible, then
run "vault init" against the target Vault instance, before encrypting and
storing the keys in the given backend.

It will not unseal the Vault instance after initializing.`,
	Run: func(cmd *cobra.Command, _ []string) {
		ctx := cmd.Context()
		slog.Info("Init Vault...")
		client, err := vaultClient.New(
			vaultClient.WithAddress("http://127.0.0.1:8200"), // TODO: Read from configuration
		)
		if err != nil {
			slog.Error(fmt.Sprintf("error connecting to Vault: %s", err.Error()))
		}
		config := customVault.Config{
			Init: customVault.InitConfig{
				KeyShares: 5,
				Threshold: 3,
			},
		}

		v, err := customVault.New(client, config)
		if err != nil {
			slog.Error(fmt.Sprintf("error creating Vault connect: %s", err.Error()))
		}

		if err = v.Init(ctx); err != nil {
			slog.Error(fmt.Sprintf("error initializing Vault: %s", err.Error()))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
