package vaultauth

import (
	"testing"

	"github.com/davepgreene/go-db-credential-refresh/store/vault/vaulttest"
)

func TestTokenAuth(t *testing.T) {
	ln, client := vaulttest.CreateTestVault(t, nil)
	defer ln.Close()

	ta := NewTokenAuth(client.Token())
	token, err := ta.GetToken(client)
	if err != nil {
		t.Error(err)
	}

	if token != client.Token() {
		t.Errorf("expected token to be %s but got %s instead", client.Token(), token)
	}
}

func TestTokenAuthWithInvalidToken(t *testing.T) {
	ln, client := vaulttest.CreateTestVault(t, nil)
	defer ln.Close()

	ta := NewTokenAuth("foobar")
	if _, err := ta.GetToken(client); err == nil {
		t.Error("expected an error but didn't get one")
	}
}
