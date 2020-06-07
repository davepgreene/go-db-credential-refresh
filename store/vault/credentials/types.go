package vaultcredentials

import (
	"github.com/davepgreene/go-db-credential-refresh/store"
	"github.com/hashicorp/vault/api"
)

// CredentialLocation represents a location where credentials can be retrieved from
type CredentialLocation interface {
	GetCredentials(client *api.Client) (string, error)
	Map(s string) (*store.Credential, error)
}

// Credentials represents an abstraction over a username and password
type Credentials interface {
	GetUsername() string
	GetPassword() string
}
