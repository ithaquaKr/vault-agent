package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the target Vault instance",
	Long: `This command will verify the backend service is accessible, then
run "vault init" against the target Vault instance, before encrypting and
storing the keys in the given backend.

It will not unseal the Vault instance after initializing.`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Print("Init Vault...")
	},
}
