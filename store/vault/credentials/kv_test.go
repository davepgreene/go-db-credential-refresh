package vaultcredentials

import (
	"context"
	"testing"

	"github.com/davepgreene/go-db-credential-refresh/store/vault/vaulttest"
)

func TestNewKvCredentials(t *testing.T) {
	ln, client := vaulttest.CreateTestVault(t)
	defer ln.Close()

	ctx := context.Background()

	path := "secret/test"
	username := "foo"
	password := "bar"

	if _, err := client.Logical().WriteWithContext(ctx, path, map[string]interface{}{
		"username": username,
		"password": password,
	}); err != nil {
		t.Fatal(err)
	}

	kvc := NewKvCredentials(path)
	credStr, err := kvc.GetCredentials(ctx, client)
	if err != nil {
		t.Fatal(err)
	}

	creds, err := kvc.Map(credStr)
	if err != nil {
		t.Fatal(err)
	}

	if creds.GetUsername() != username {
		t.Fatalf("expected username to be %s but got %s", username, creds.GetUsername())
	}

	if creds.GetPassword() != password {
		t.Fatalf("expected password to be %s but got %s instead", password, creds.GetPassword())
	}
}
