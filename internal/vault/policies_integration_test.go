//go:build integration

package vault

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	vaultClient "github.com/hashicorp/vault-client-go"
	"github.com/ithaquaKr/vault-agent/pkg/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Integration Test for Vault
type VaultIntegrationTestSuite struct {
	suite.Suite
	vaultContainer *helper.VaultContainer
	vault          *vaultServer
	ctx            context.Context
}

func (suite *VaultIntegrationTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	vaultContainer, err := helper.CreateVaultContainer(suite.ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("error while creating test container: %s", err.Error()))
	}
	suite.vaultContainer = vaultContainer
	client, err := vaultClient.New(vaultClient.WithAddress(vaultContainer.Address))
	if err != nil {
		slog.Error(fmt.Sprintf("error connecting to Vault: %s", err.Error()))
	}
	config := Config{}
	suite.vault = &vaultServer{
		cl:     client,
		config: &config,
	}
}

func (suite *VaultIntegrationTestSuite) TearDownSuite() {
	if err := suite.vaultContainer.Terminate(suite.ctx); err != nil {
		slog.Error(fmt.Sprintf("error terminating vault container: %s", err.Error()))
	}
}

// Test Policies
func (suite *VaultIntegrationTestSuite) TestGetExistingPolicies() {
	t := suite.T()
	_, err := suite.vault.getExistingPolicies()
	assert.NoError(t, err)
}

func TestVaultIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(VaultIntegrationTestSuite))
}
