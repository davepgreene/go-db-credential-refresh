package vaultcredentials

import (
	"testing"

	"github.com/davepgreene/go-db-credential-refresh/store/vault/vaulttest"
)

func TestNewKvCredentials(t *testing.T) {
	ln, client := vaulttest.CreateTestVault(t, nil)
	defer ln.Close()

	path := "secret/test"
	username := "foo"
	password := "bar"

	if _, err := client.Logical().Write(path, map[string]interface{}{
		"username": username,
		"password": password,
	}); err != nil {
		t.Error(err)
	}

	kvc := NewKvCredentials(path)
	credStr, err := kvc.GetCredentials(client)
	if err != nil {
		t.Error(err)
	}

	creds, err := kvc.Map(credStr)
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
