package vault

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/api"

	"github.com/davepgreene/go-db-credential-refresh/driver"
	vaultauth "github.com/davepgreene/go-db-credential-refresh/store/vault/auth"
	vaultcredentials "github.com/davepgreene/go-db-credential-refresh/store/vault/credentials"
)

type TokenLocation interface {
	GetToken(client *api.Client) (string, error)
}

// Store is a Store implementation for HashiCorp Vault.
type Store struct {
	client *api.Client
	cl     vaultcredentials.CredentialLocation
	tl     TokenLocation
	creds  driver.Credentials
}

// Config contains configuration information.
type Config struct {
	Client             *api.Client
	TokenLocation      TokenLocation
	CredentialLocation vaultcredentials.CredentialLocation
}

var (
	ErrConfigRequired             = errors.New("config is required")
	ErrCredentialLocationRequired = errors.New("credential location is required")
	ErrClientRequired             = errors.New("client is required")
	ErrTokenLocationRequired      = errors.New("token location is required")
)

// NewStore creates a new Vault-backed store.
func NewStore(c *Config) (*Store, error) {
	if c == nil {
		return nil, ErrConfigRequired
	}

	if c.CredentialLocation == nil {
		return nil, ErrCredentialLocationRequired
	}

	client := c.Client
	if client == nil {
		return nil, ErrClientRequired
	}

	if c.TokenLocation == nil {
		// If the token location is nil, we should check if the client already has a token
		if client.Token() == "" {
			return nil, ErrTokenLocationRequired
		}

		c.TokenLocation = vaultauth.NewTokenAuth(client.Token())
	}

	token, err := c.TokenLocation.GetToken(client)
	if err != nil {
		return nil, err
	}

	client.SetToken(token)

	return &Store{
		client: client,
		tl:     c.TokenLocation,
		cl:     c.CredentialLocation,
	}, nil
}

// Get implements the Store interface.
func (v *Store) Get(ctx context.Context) (driver.Credentials, error) {
	if v.creds != nil {
		return v.creds, nil
	}

	return v.Refresh(ctx)
}

// Refresh implements the store interface.
func (v *Store) Refresh(ctx context.Context) (driver.Credentials, error) {
	credStr, err := v.cl.GetCredentials(v.client)
	if err != nil {
		return nil, err
	}

	creds, err := v.cl.Map(credStr)
	if err != nil {
		return nil, err
	}

	// Cache the credentials
	v.creds = creds

	return creds, nil
}
