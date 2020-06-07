package vaulttest

import (
	"net"
	"testing"

	log "github.com/hashicorp/go-hclog"
	credKube "github.com/hashicorp/vault-plugin-auth-kubernetes"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/database"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func CreateTestVault(t *testing.T, l log.Logger) (net.Listener, *api.Client) {
	t.Helper()

	coreConf := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"kubernetes": credKube.Factory,
		},
		LogicalBackends: map[string]logical.Factory{
			"database": database.Factory,
		},
	}

	if l != nil {
		coreConf.Logger = l
	}

	core, keyShares, rootToken := vault.TestCoreUnsealedWithConfig(t, coreConf)

	_ = keyShares

	// Start an HTTP server for the core.
	ln, addr := http.TestServer(t, core)

	// Create a client that talks to the server, initially authenticating with
	// the root token.
	conf := api.DefaultConfig()
	conf.Address = addr

	client, err := api.NewClient(conf)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(rootToken)

	return ln, client
}
