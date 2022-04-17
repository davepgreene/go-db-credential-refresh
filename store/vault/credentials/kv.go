package vaultcredentials

import (
	"context"

	"github.com/hashicorp/vault/api"

	"github.com/davepgreene/go-db-credential-refresh/store"
)

// KvCredentials implements the CredentialLocation interface.
type KvCredentials struct {
	path string
}

// NewKvCredentials retrieves credentials from Vault's K/V store.
func NewKvCredentials(path string) CredentialLocation {
	return &KvCredentials{
		path: path,
	}
}

// GetCredentials implements the CredentialLocation interface.
func (kv *KvCredentials) GetCredentials(ctx context.Context, client *api.Client) (string, error) {
	return GetFromVaultSecretsAPI(ctx, client, kv.path)
}

// Map implements the CredentialLocation interface.
func (*KvCredentials) Map(s string) (*store.Credential, error) {
	return DefaultMapper(s)
}
