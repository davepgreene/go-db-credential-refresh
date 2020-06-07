package vaultcredentials

import (
	"encoding/json"
	"testing"

	"github.com/davepgreene/go-db-credential-refresh/store/vault/vaulttest"
)

func TestGetFromVaultSecretsAPI(t *testing.T) {
	ln, client := vaulttest.CreateTestVault(t, nil)
	defer ln.Close()

	// Valid path with response
	b, err := GetFromVaultSecretsAPI(client, "auth/token/lookup-self")
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
	_, err = GetFromVaultSecretsAPI(client, "flerp/derp/herp")
	if err == nil {
		t.Error("expected an error but didn't get one")
	}

	if err != errInvalidPath {
		t.Errorf("expected a '%T' but got '%T' instead", errInvalidPath, err)
	}
}

func TestGetFromVaultSecretsAPIWithVaultError(t *testing.T) {
	ln, client := vaulttest.CreateTestVault(t, nil)
	defer ln.Close()

	if _, err := client.Logical().Write("secret/foo", map[string]interface{}{
		"secret": "string",
	}); err != nil {
		t.Error(err)
	}

	if _, err := client.Logical().Write("sys/policy/restricted", map[string]interface{}{
		"policy": `path "secret/foo" {
			capabilities = ["deny"]
		}`,
	}); err != nil {
		t.Error(err)
	}

	resp, err := client.Logical().Write("auth/token/create", map[string]interface{}{
		"policies": []string{"restricted"},
	})
	if err != nil {
		t.Error(err)
	}

	client.SetToken(resp.Auth.ClientToken)

	if resp, err := GetFromVaultSecretsAPI(client, "secret/foo"); err == nil {
		t.Errorf("expected an error but got '%s' as a response instead", resp)
	}
}
