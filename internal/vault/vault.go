package vault

import (
	"errors"
	"fmt"

	vaultApi "github.com/hashicorp/vault/api"
)

type VaultInitConfig struct {
	KeyShares int `mapstructure:"keyShares"`
	Threshold int `mapstructure:"threshold"`
}

type VaultData struct {
	Policies []policy `mapstructure:"policies"`
}

type vaultController struct {
	cl         *vaultApi.Client
	initConfig *VaultInitConfig
	data       *VaultData
}

// New create a vaultController to do action with Vault
func New(cl *vaultApi.Client, initConfig VaultInitConfig, data VaultData) (*vaultController, error) {
	return &vaultController{
		cl:         cl,
		initConfig: &initConfig,
		data:       &data,
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
