package vaultcredentials

import (
	"fmt"
	"testing"

	"github.com/davepgreene/go-db-credential-refresh/store/vault/vaulttest"
)

func TestNewAPIDatabaseCredentials(t *testing.T) {
	ln, client := vaulttest.CreateTestVault(t, nil)
	defer ln.Close()

	// Because this CredentialLocation is agnostic to the location of the actual credentials we can
	// fudge this test by using the k/v secret type rather than building a mock vault plugin,
	// mounting it as a db type, and dealing with vault's complicated "separate binary with gRPC
	// communication" process.
	// Instead we mount the k/v secret type at `database`.
	if _, err := client.Logical().Write("sys/mounts/database", map[string]interface{}{
		"type": "kv",
	}); err != nil {
		t.Error(err)
	}

	role := "postgres"

	if _, err := client.Logical().Write(fmt.Sprintf("database/creds/%s", role), map[string]interface{}{
		"username": username,
		"password": password,
	}); err != nil {
		t.Error(err)
	}

	adc := NewAPIDatabaseCredentials(role, "")
	credStr, err := adc.GetCredentials(client)
	if err != nil {
		t.Error(err)
	}

	creds, err := adc.Map(credStr)
	if err != nil {
		t.Error(err)
	}

	if creds.GetUsername() != username {
		t.Errorf("expected username to be %s but got %s", username, creds.GetUsername())
	}

	if creds.GetPassword() != password {
		t.Errorf("expected password to be %s but got %s instead", password, creds.GetPassword())
	}
}
