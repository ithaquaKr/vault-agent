package vault

import (
	"fmt"
	"log/slog"

	vaultApi "github.com/hashicorp/vault/api"
)

// Init initializes the vault instance
func (v *vaultManager) Init() error {
	isInitialized, err := v.cl.Sys().InitStatus()
	if err != nil {
		return fmt.Errorf("error checking vault initialized status: %w", err)
	}
	if isInitialized {
		slog.Info("vault is already initialized")
		return nil
	}

	slog.Info("initializing vault...")
	initReq := vaultApi.InitRequest{
		SecretShares:    v.initConfig.KeyShares,
		SecretThreshold: v.initConfig.Threshold,
	}
	resp, err := v.cl.Sys().Init(&initReq)
	if err != nil {
		return fmt.Errorf("error initializing vault: %w", err)
	}

	for i, k := range resp.Keys {
		slog.Info(fmt.Sprintf("Unseal key %s: %s", keyUnsealForID(i), k))
	}

	slog.Info(fmt.Sprintf("Token root: %s", resp.RootToken))

	return nil
}

// Unseal unlocks the vault using the unseal keys.
func (v *vaultManager) Unseal() error {
	return nil
}

func keyUnsealForID(i int) string {
	return fmt.Sprintf("vault-unseal-%d", i+1)
}
