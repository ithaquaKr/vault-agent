package vault

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	vaultClient "github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

type InitConfig struct {
	KeyShares int32
	Threshold int32
}

type Config struct {
	Init     InitConfig `mapstructure:"InitConfig"`
	Policies []policy   `mapstructure:"policies"`
}

type vaultServer struct {
	ctx    context.Context
	cl     *vaultClient.Client
	config *Config
}

// New returns a new implement of Operator interface, or an error
func New(cl *vaultClient.Client, config Config, ctx context.Context) (*vaultServer, error) {
	return &vaultServer{
		cl:     cl,
		config: &config,
		ctx:    ctx,
	}, nil
}

// IsSealed determine if Vault is sealed.
func (v *vaultServer) IsSealed(ctx context.Context) (bool, error) {
	resp, err := v.cl.System.SealStatus(ctx)
	if err != nil {
		return false, errors.New("error checking status")
	}
	return resp.Data.Sealed, nil
}

// Unseal the vault instance
func (v *vaultServer) Unseal() error {
	return nil
}

// Leader check if instance is Leader.
func (v *vaultServer) Leader(ctx context.Context) (bool, error) {
	resp, err := v.cl.System.LeaderStatus(ctx)
	if err != nil {
		return false, errors.New("error checking leader")
	}
	return resp.Data.IsSelf, nil
}

// LeaderAddress check leader address.
func (v *vaultServer) LeaderAddress(ctx context.Context) (string, error) {
	resp, err := v.cl.System.LeaderStatus(ctx)
	if err != nil {
		return "", fmt.Errorf("error checking leader, err: %s", err)
	}

	return resp.Data.LeaderAddress, nil
}

func (v *vaultServer) Init(ctx context.Context) error {
	isInitializedResp, err := v.cl.System.ReadInitializationStatus(ctx)
	if err != nil {
		return fmt.Errorf("error checking vault initialized status: %s", err.Error())
	}
	if isInitializedResp.Data["initialized"] == true {
		slog.Info("vault is already initialized")
		return nil
	}

	slog.Info("initializing vault...")
	initReq := schema.InitializeRequest{
		SecretShares:    v.config.Init.KeyShares,
		SecretThreshold: v.config.Init.Threshold,
	}
	resp, err := v.cl.System.Initialize(ctx, initReq)
	if err != nil {
		return fmt.Errorf("error initializing vault: %s", err.Error())
	}

	if keys, ok := resp.Data["keys"].([]interface{}); ok {
		for i, k := range keys {
			slog.Info(fmt.Sprintf("Unseal key %s: %s", keyUnsealForID(i), k))
		}
	}

	return nil
}

func keyUnsealForID(i int) string {
	return fmt.Sprintf("vault-unseal-%d", i+1)
}
