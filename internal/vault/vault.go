package vault

import (
	"errors"
	"fmt"

	vaultApi "github.com/hashicorp/vault/api"

	"github.com/ithaquaKr/vault-agent/pkg/config"
)

type vaultManager struct {
	cl         *vaultApi.Client
	initConfig *config.VaultInitConfig
	data       *config.VaultData
}

// New create a vaultController to do action with Vault
func New(cl *vaultApi.Client, initConfig config.VaultInitConfig, data config.VaultData) (*vaultManager, error) {
	return &vaultManager{
		cl:         cl,
		initConfig: &initConfig,
		data:       &data,
	}, nil
}

// IsSealed determine if Vault is sealed.
func (v *vaultManager) IsSealed() (bool, error) {
	resp, err := v.cl.Sys().SealStatus()
	if err != nil {
		return false, errors.New("error checking status")
	}
	return resp.Sealed, nil
}

// Leader check if instance is Leader.
func (v *vaultManager) Leader() (bool, error) {
	resp, err := v.cl.Sys().Leader()
	if err != nil {
		return false, errors.New("error checking leader")
	}
	return resp.IsSelf, nil
}

// LeaderAddress check leader address.
func (v *vaultManager) LeaderAddress() (string, error) {
	resp, err := v.cl.Sys().Leader()
	if err != nil {
		return "", fmt.Errorf("error checking leader, err: %s", err)
	}

	return resp.LeaderAddress, nil
}
