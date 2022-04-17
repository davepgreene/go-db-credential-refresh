package vaultcredentials

import (
	"context"
	"fmt"
	"testing"

	"github.com/davepgreene/go-db-credential-refresh/store/vault/vaulttest"
)

func TestNewAPIDatabaseCredentials(t *testing.T) {
	ln, client := vaulttest.CreateTestVault(t, nil)
	defer ln.Close()

	ctx := context.Background()

	// Because this CredentialLocation is agnostic to the location of the actual credentials we can
	// fudge this test by using the k/v secret type rather than building a mock vault plugin,
	// mounting it as a db type, and dealing with vault's complicated "separate binary with gRPC
	// communication" process.
	// Instead we mount the k/v secret type at `database`.
	if _, err := client.Logical().WriteWithContext(ctx, "sys/mounts/database", map[string]interface{}{
		"type": "kv",
	}); err != nil {
		t.Error(err)
	}

	role := "postgres"

	if _, err := client.Logical().WriteWithContext(ctx, fmt.Sprintf("database/creds/%s", role), map[string]interface{}{
		"username": username,
		"password": password,
	}); err != nil {
		t.Error(err)
	}

	adc := NewAPIDatabaseCredentials(role, "")
	credStr, err := adc.GetCredentials(ctx, client)
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
