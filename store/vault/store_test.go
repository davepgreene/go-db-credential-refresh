package vault

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/vault/api"

	"github.com/davepgreene/go-db-credential-refresh/store"
	vaultcredentials "github.com/davepgreene/go-db-credential-refresh/store/vault/credentials"
	"github.com/davepgreene/go-db-credential-refresh/store/vault/vaulttest"
)

const (
	token    = "token"
	username = "foo"
	password = "bar"
)

type testTokenLocation struct {
	TokenGetter func(ctx context.Context, c *api.Client) (string, error)
}

func (k *testTokenLocation) GetToken(ctx context.Context, client *api.Client) (string, error) {
	return k.TokenGetter(ctx, client)
}

type testCredentialLocation struct {
	CredentialGetter func(ctx context.Context, c *api.Client) (string, error)
	Mapper           func(s string) (*store.Credential, error)
}

func (tcl *testCredentialLocation) GetCredentials(ctx context.Context, client *api.Client) (string, error) {
	return tcl.CredentialGetter(ctx, client)
}

func (tcl *testCredentialLocation) Map(s string) (*store.Credential, error) {
	return tcl.Mapper(s)
}

func TestNewStoreCannotCreateWithoutValidConfig(t *testing.T) {
	if _, err := NewStore(nil); err == nil {
		t.Fatal("expected an error but didn't get one")
	}

	if _, err := NewStore(&Config{}); err == nil {
		t.Fatal("expected an error but didn't get one")
	}

	client, err := api.NewClient(&api.Config{
		Address: "localhost:8200",
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := NewStore(&Config{
		Client: client,
	}); err == nil {
		t.Fatal("expected an error but didn't get one")
	}

	if _, err := NewStore(&Config{
		Client:        client,
		TokenLocation: &testTokenLocation{},
	}); err == nil {
		t.Fatal("expected an error but didn't get one")
	}

	if _, err := NewStore(&Config{
		TokenLocation:      &testTokenLocation{},
		CredentialLocation: &testCredentialLocation{},
	}); err == nil {
		t.Fatal("expected an error but didn't get one")
	}

	if _, err := NewStore(&Config{
		Client: client,
		TokenLocation: &testTokenLocation{
			TokenGetter: func(_ context.Context, _ *api.Client) (string, error) {
				return "", errors.New("unable to get token")
			},
		},
		CredentialLocation: &testCredentialLocation{},
	}); err == nil {
		t.Fatal("expected an error but didn't get one")
	}
}

func TestNewStoreWithValidConfig(t *testing.T) {
	client, err := api.NewClient(&api.Config{
		Address: "localhost:8200",
	})
	if err != nil {
		t.Fatal(err)
	}

	s, err := NewStore(&Config{
		Client: client,
		TokenLocation: &testTokenLocation{
			TokenGetter: func(_ context.Context, _ *api.Client) (string, error) {
				return token, nil
			},
		},
		CredentialLocation: &testCredentialLocation{
			CredentialGetter: func(_ context.Context, _ *api.Client) (string, error) {
				return fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password), nil
			},
			Mapper: vaultcredentials.DefaultMapper,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	creds, err := s.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if creds.GetUsername() != username {
		t.Fatalf("expected username to be '%s' but got '%s' instead", username, creds.GetUsername())
	}

	if creds.GetPassword() != password {
		t.Fatalf("expected password to be '%s' but got '%s' instead", password, creds.GetPassword())
	}
}

func TestNewStoreWithGetCredentialError(t *testing.T) {
	client, err := api.NewClient(&api.Config{
		Address: "localhost:8200",
	})
	if err != nil {
		t.Fatal(err)
	}

	s, err := NewStore(&Config{
		Client: client,
		TokenLocation: &testTokenLocation{
			TokenGetter: func(_ context.Context, _ *api.Client) (string, error) {
				return token, nil
			},
		},
		CredentialLocation: &testCredentialLocation{
			CredentialGetter: func(_ context.Context, _ *api.Client) (string, error) {
				return "", errors.New("could not get credentials")
			},
			Mapper: func(s string) (*store.Credential, error) {
				return nil, nil
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	if _, err := s.Get(ctx); err == nil {
		t.Fatal("expected an error but didn't get one")
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewStoreWithCredentialMapperError(t *testing.T) {
	client, err := api.NewClient(&api.Config{
		Address: "localhost:8200",
	})
	if err != nil {
		t.Fatal(err)
	}

	s, err := NewStore(&Config{
		Client: client,
		TokenLocation: &testTokenLocation{
			TokenGetter: func(_ context.Context, _ *api.Client) (string, error) {
				return token, nil
			},
		},
		CredentialLocation: &testCredentialLocation{
			CredentialGetter: func(_ context.Context, _ *api.Client) (string, error) {
				return fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password), nil
			},
			Mapper: func(s string) (*store.Credential, error) {
				return nil, errors.New("failed to unmarshal credentials")
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	if _, err := s.Get(ctx); err == nil {
		t.Fatal("expected an error but didn't get one")
	}
}

func TestNewStoreWithClientThatAlreadyHasToken(t *testing.T) {
	ln, client := vaulttest.CreateTestVault(t)
	defer ln.Close()

	s, err := NewStore(&Config{
		Client: client,
		CredentialLocation: &testCredentialLocation{
			CredentialGetter: func(_ context.Context, c *api.Client) (string, error) {
				if c.Token() != client.Token() {
					t.Fatalf("expected token to be '%s' but got '%s' instead", client.Token(), c.Token())
				}

				return fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password), nil
			},
			Mapper: vaultcredentials.DefaultMapper,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	creds, err := s.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if creds.GetUsername() != username {
		t.Fatalf("expected username to be '%s' but got '%s' instead", username, creds.GetUsername())
	}

	if creds.GetPassword() != password {
		t.Fatalf("expected password to be '%s' but got '%s' instead", password, creds.GetPassword())
	}
}

func TestNewStoreWithInvalidTokenLocation(t *testing.T) {
	envToken := os.Getenv(api.EnvVaultToken)
	client, err := api.NewClient(&api.Config{
		Address: "localhost:8200",
	})
	if err != nil {
		t.Fatal(err)
	}

	// If the client has pulled a token from the environment we deliberately unset it to mimic
	// a scenario where there's no token present in any way.
	if client.Token() == envToken {
		client.SetToken("")
	}

	if _, err := NewStore(&Config{
		Client: client,
		CredentialLocation: &testCredentialLocation{
			CredentialGetter: func(_ context.Context, _ *api.Client) (string, error) {
				return "", nil
			},
			Mapper: func(s string) (*store.Credential, error) {
				return nil, nil
			},
		},
	}); err == nil {
		t.Fatal("expected an error but didn't get one")
	}
}

func TestStoreWithCachedCredentials(t *testing.T) {
	ln, client := vaulttest.CreateTestVault(t)
	defer ln.Close()

	mapCallCount := 0

	s, err := NewStore(&Config{
		Client: client,
		CredentialLocation: &testCredentialLocation{
			CredentialGetter: func(_ context.Context, _ *api.Client) (string, error) {
				return fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password), nil
			},
			Mapper: func(s string) (*store.Credential, error) {
				mapCallCount++

				return vaultcredentials.DefaultMapper(s)
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	creds, err := s.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if creds.GetUsername() != username {
		t.Fatalf("expected username to be '%s' but got '%s' instead", username, creds.GetUsername())
	}

	if creds.GetPassword() != password {
		t.Fatalf("expected password to be '%s' but got '%s' instead", password, creds.GetPassword())
	}

	if _, err = s.Get(ctx); err != nil {
		t.Fatal(err)
	}

	if mapCallCount != 1 {
		t.Fatalf("expected the mapper function to only be called once but it was called %d times", mapCallCount)
	}
}
