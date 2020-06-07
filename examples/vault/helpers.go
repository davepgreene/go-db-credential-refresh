package main

import (
	"errors"
	"fmt"
	"net"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"

	"github.com/davepgreene/go-db-credential-refresh/store/vault/vaulttest"
)

// SetupVault sets up an in-memory Vault core and webserver then enables the plugins/configs we need for
// this example.
func SetupVault() (net.Listener, *api.Client, error) {
	fmt.Println("Creating in-memory Vault instance")

	t := &testing.T{}
	ln, client := vaulttest.CreateTestVault(t, log.NewNullLogger())

	fmt.Println("Mounting the database backend")
	if _, err := client.Logical().Write("sys/mounts/database", map[string]interface{}{
		"type": "database",
	}); err != nil {
		return nil, nil, err
	}

	uri := fmt.Sprintf("postgresql://{{username}}:{{password}}@%s:%d/?sslmode=disable", host, port)

	fmt.Println("Configuring the postgres database and role")
	if _, err := client.Logical().Write(fmt.Sprintf("database/config/%s", dbName), map[string]interface{}{
		"plugin_name":    "postgresql-database-plugin",
		"allowed_roles":  role,
		"connection_url": uri,
		"username":       username,
		"password":       password,
	}); err != nil {
		return nil, nil, err
	}

	if _, err := client.Logical().Write(fmt.Sprintf("database/roles/%s", role), map[string]interface{}{
		"db_name": dbName,
		"creation_statements": []string{
			`CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'`,
			`GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{name}}"`,
		},
		"default_ttl": "2s",
		"max_ttl":     "5s",
	}); err != nil {
		return nil, nil, err
	}

	fmt.Println("Vault has been configured")

	return ln, client, nil
}

// TearDownRoles invalidates any Vault leases which will delete those roles in the DB
func TearDownRoles(client *api.Client) error {
	fmt.Println("Removing existing roles from Vault before shutdown")

	pathTemplate := "sys/leases/%s/database/creds/%s"
	resp, err := client.Logical().List(fmt.Sprintf(pathTemplate, "lookup", role))
	if err != nil {
		return err
	}
	if resp == nil {
		return errors.New("invalid path for lease lookup")
	}

	if keys, ok := resp.Data["keys"]; ok {
		leases := keys.([]interface{})
		fmt.Println("") // so there's a line between the ctrl-c character
		for _, l := range leases {
			lease := l.(string)
			fmt.Printf("Revoking lease %s...\n", lease)
			leasePath := fmt.Sprintf("%s/%s", fmt.Sprintf(pathTemplate, "revoke", role), lease)
			if _, err = client.Logical().Write(leasePath, nil); err != nil {
				return err
			}
			fmt.Printf("Successfully revoked lease %s...\n", lease)
		}
	}

	fmt.Println("Successfully removed roles from Vault which dropped created DB users")

	return nil
}
