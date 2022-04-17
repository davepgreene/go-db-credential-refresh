package vaultcredentials

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/davepgreene/go-db-credential-refresh/store/vault/vaulttest"
)

func TestGetFromVaultSecretsAPI(t *testing.T) {
	ln, client := vaulttest.CreateTestVault(t, nil)
	defer ln.Close()

	ctx := context.Background()

	// Valid path with response
	b, err := GetFromVaultSecretsAPI(ctx, client, "auth/token/lookup-self")
	if err != nil {
		t.Fatal(err)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(b), &resp); err != nil {
		t.Error(err)
	}
	// Testing secret data attributes is a bit brittle unfortunately :(
	for _, field := range []string{
		"accessor",
		"creation_time",
		"creation_ttl",
		"id",
		"path",
		"ttl",
		"type",
	} {
		if _, ok := resp[field]; !ok {
			t.Errorf("expected '%s' to be in response data", field)
		}
	}

	// Invalid path
	_, err = GetFromVaultSecretsAPI(ctx, client, "flerp/derp/herp")
	if err == nil {
		t.Error("expected an error but didn't get one")
	}

	if !errors.Is(err, errInvalidPath) {
		t.Errorf("expected a '%T' but got '%T' instead", errInvalidPath, err)
	}
}

func TestGetFromVaultSecretsAPIWithVaultError(t *testing.T) {
	ln, client := vaulttest.CreateTestVault(t, nil)
	defer ln.Close()

	ctx := context.Background()

	if _, err := client.Logical().WriteWithContext(ctx, "secret/foo", map[string]interface{}{
		"secret": "string",
	}); err != nil {
		t.Error(err)
	}

	if _, err := client.Logical().WriteWithContext(ctx, "sys/policy/restricted", map[string]interface{}{
		"policy": `path "secret/foo" {
			capabilities = ["deny"]
		}`,
	}); err != nil {
		t.Error(err)
	}

	resp, err := client.Logical().WriteWithContext(ctx, "auth/token/create", map[string]interface{}{
		"policies": []string{"restricted"},
	})
	if err != nil {
		t.Error(err)
	}

	client.SetToken(resp.Auth.ClientToken)

	if resp, err := GetFromVaultSecretsAPI(ctx, client, "secret/foo"); err == nil {
		t.Errorf("expected an error but got '%s' as a response instead", resp)
	}
}
