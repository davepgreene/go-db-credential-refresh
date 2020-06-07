package vaultcredentials

import (
	"fmt"

	"github.com/davepgreene/go-db-credential-refresh/store"
	"github.com/hashicorp/vault/api"
)

// APIDatabaseCredentials gets DB credentials from the Vault Database Secrets engine
// See: https://www.vaultproject.io/docs/secrets/databases/
type APIDatabaseCredentials struct {
	path string
	role string
}

// NewAPIDatabaseCredentials creates a new credential location backed by Vault's DB Secrets engine.
//
// The path argument will be mostly unused unless the user mounts the database backend in a different
// location.
func NewAPIDatabaseCredentials(role, path string) CredentialLocation {
	if path == "" {
		path = "database"
	}
	return &APIDatabaseCredentials{
		role: role,
		path: path,
	}
}

// GetCredentials implements the CredentialLocation interface
func (db *APIDatabaseCredentials) GetCredentials(client *api.Client) (string, error) {
	return GetFromVaultSecretsAPI(client, fmt.Sprintf("%s/creds/%s", db.path, db.role))
}

// Map implements the CredentialLocation interface
func (db *APIDatabaseCredentials) Map(s string) (*store.Credential, error) {
	return DefaultMapper(s)
}
