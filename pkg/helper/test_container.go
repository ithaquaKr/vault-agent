package helper

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go/modules/vault"
)

type VaultContainer struct {
	*vault.VaultContainer
	Address string
}

func CreateVaultContainer(ctx context.Context) (*VaultContainer, error) {
	vaultContainer, err := vault.Run(ctx, "hashicorp/vault:1.16.1")
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %s", err.Error())
	}

	vaultAddress, err := vaultContainer.Host(ctx)

	return &VaultContainer{
		VaultContainer: vaultContainer,
		Address:        vaultAddress,
	}, nil
}
