package vaultcredentials

import (
	"context"

	"github.com/hashicorp/vault/api"

	"github.com/davepgreene/go-db-credential-refresh/store"
)

// CredentialLocation represents a location where credentials can be retrieved from.
type CredentialLocation interface {
	GetCredentials(ctx context.Context, client *api.Client) (string, error)
	Map(s string) (*store.Credential, error)
}

// Credentials represents an abstraction over a username and password.
type Credentials interface {
	GetUsername() string
	GetPassword() string
}
