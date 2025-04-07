package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var unsealCmd = &cobra.Command{
	Use:   "unseal",
	Short: "Unseal the target Vault instance",
	Long: `Unseals Vault with unseal keys provide from command line.

  It will continuously attempt to unseal the target Vault instance, 
  by retrieving unseal keys from command line.`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Print("Unsealing Vault..")
	},
}
