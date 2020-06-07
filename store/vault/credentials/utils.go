package vaultcredentials

import (
	"encoding/json"
	"errors"

	"github.com/hashicorp/vault/api"
)

var (
	errInvalidPath = errors.New("invalid path")
)

// GetFromVaultSecretsAPI is a wrapper over logical reads from a Vault path with marshalling and error handling.
func GetFromVaultSecretsAPI(client *api.Client, path string) (string, error) {
	resp, err := client.Logical().Read(path)
	if err != nil {
		return "", err
	}

	// If Vault can't handle the path it will return a nil response with no error
	// so it's important to nil check it so we don't accidentally try to marshal it.
	if resp == nil {
		return "", errInvalidPath
	}

	// Something in Vault's API would have to be horribly broken for the response
	// not to be marshalable but it's worth error checking it as a matter of habit.
	b, err := json.Marshal(resp.Data)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
