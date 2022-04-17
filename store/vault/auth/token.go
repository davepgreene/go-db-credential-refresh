package vaultauth

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/api"
)

const (
	tokenSelfLookupPath = "auth/token/lookup-self" //nolint:gosec
)

// TokenAuth is a pass-through authentication mechanism to set vault tokens directly for
// use by the Vault store.
// NOTE: Token renewal should be handled outside of this library.
type TokenAuth struct {
	token string
}

var (
	ErrUnableToLookupToken = errors.New("unable to lookup token information")
)

// NewTokenAuth creates a new Vault token auth location.
func NewTokenAuth(token string) *TokenAuth {
	return &TokenAuth{
		token: token,
	}
}

// GetToken implements the TokenLocation interface.
func (t *TokenAuth) GetToken(ctx context.Context, client *api.Client) (string, error) {
	client.SetToken(t.token)
	// Before we pass the token back we should call an endpoint it will have access to just to be sure
	resp, err := client.Logical().ReadWithContext(ctx, tokenSelfLookupPath)
	if err != nil {
		return "", err
	}
	// We could hit this branch if Vault's token `lookup-self` path is removed but that's pretty unlikely
	// to happen and if it does I'm sure many other things will have broken well before then.
	if resp == nil {
		return "", ErrUnableToLookupToken
	}

	return t.token, nil
}
