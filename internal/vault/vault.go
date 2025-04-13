package vault

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime"

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
func (v *vaultController) Unseal(keys []string, addresses []string) error {
	defer runtime.GC()

	// Infinity loop for unseal Vault
	for i := 0; ; i++ {
		slog.Info("Start unsealing Vault...")
		for _, address := range addresses {
			slog.Debug(fmt.Sprintf("unsealing instance with address: %s", address))

			keyNum := len(keys)
			configKeyNum := v.config.Init.Threshold
			if keyNum != configKeyNum {
				return fmt.Errorf("number of keys is not equal threshold, %d not equal %d", keyNum, configKeyNum)
			}

			for _, key := range keys {
				resp, err := v.cl.Sys().Unseal(string(key))
				if err != nil {
					return fmt.Errorf("fail to send unseal request to vault, err: %s", err)
				}

				if !resp.Sealed {
					return nil
				}

				if resp.Progress == 0 {
					return fmt.Errorf("fail to unseal vault")
				}
			}
		}
	}
}

// keyUnsealForID [TODO:description]
func keyUnsealForID(i int) string {
	return fmt.Sprintf("vault-unseal-%d", i+1)
}
