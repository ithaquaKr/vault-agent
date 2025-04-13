package client

import (
	"net/http"

	vaultApi "github.com/hashicorp/vault/api"
)

// TODO: Optimize this

// NewVaultClient create a new client to interact with Vault
func NewVaultClient(addr string, insecureSkipVerify bool) (*vaultApi.Client, error) {
	cfg := vaultApi.DefaultConfig()
	cfg.Address = addr

	clientTLSConfig := cfg.HttpClient.Transport.(*http.Transport).TLSClientConfig
	if insecureSkipVerify {
		clientTLSConfig.InsecureSkipVerify = true
	}

	return vaultApi.NewClient(cfg)
}
