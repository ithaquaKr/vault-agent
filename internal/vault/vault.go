package vault

import (
	"errors"
	"fmt"
	"log/slog"

	vaultApi "github.com/hashicorp/vault/api"
)

type InitConfig struct {
	KeyShares int `mapstructure:"keyShares"`
	Threshold int `mapstructure:"threshold"`
}

type VaultConfig struct {
	Init     InitConfig `mapstructure:"initConfig"`
	Policies []policy   `mapstructure:"policies"`
}

type vaultController struct {
	cl     *vaultApi.Client
	config *VaultConfig
}

// New create a vaultManager to do action with Vault
func New(cl *vaultApi.Client, config VaultConfig) (*vaultController, error) {
	return &vaultController{
		cl:     cl,
		config: &config,
	}, nil
}

// IsSealed determine if Vault is sealed.
func (v *vaultController) IsSealed() (bool, error) {
	resp, err := v.cl.Sys().SealStatus()
	if err != nil {
		return false, errors.New("error checking status")
	}
	return resp.Sealed, nil
}

// Leader check if instance is Leader.
func (v *vaultController) Leader() (bool, error) {
	resp, err := v.cl.Sys().Leader()
	if err != nil {
		return false, errors.New("error checking leader")
	}
	return resp.IsSelf, nil
}

// LeaderAddress check leader address.
func (v *vaultController) LeaderAddress() (string, error) {
	resp, err := v.cl.Sys().Leader()
	if err != nil {
		return "", fmt.Errorf("error checking leader, err: %s", err)
	}

	return resp.LeaderAddress, nil
}

func (v *vaultController) Init() error {
	isInitialized, err := v.cl.Sys().InitStatus()
	if err != nil {
		return fmt.Errorf("error checking vault initialized status: %s", err.Error())
	}
	if isInitialized {
		slog.Info("vault is already initialized")
		return nil
	}

	slog.Info("initializing vault...")
	initReq := vaultApi.InitRequest{
		SecretShares:    v.config.Init.KeyShares,
		SecretThreshold: v.config.Init.Threshold,
	}
	resp, err := v.cl.Sys().Init(&initReq)
	if err != nil {
		return fmt.Errorf("error initializing vault: %s", err.Error())
	}

	for i, k := range resp.Keys {
		slog.Info(fmt.Sprintf("Unseal key %s: %s", keyUnsealForID(i), k))
	}

	slog.Info(fmt.Sprintf("Token root: %s", resp.RootToken))

	return nil
}

// Unseal the vault instance
func (v *vaultController) Unseal() error {
	return nil
}

func keyUnsealForID(i int) string {
	return fmt.Sprintf("vault-unseal-%d", i+1)
}
