package vaultauth

import (
	"io/ioutil"

	"github.com/hashicorp/vault/api"
)

const (
	kubeTokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token" //nolint:gosec
)

// KubernetesAuth gets the vault auth token from the kubernetes secrets file.
type KubernetesAuth struct {
	role string
	path string
}

// NewKubernetesAuth creates a new k8s secret auth token location.
func NewKubernetesAuth(role, path string) *KubernetesAuth {
	if path == "" {
		path = kubeTokenPath
	}

	return &KubernetesAuth{
		role: role,
		path: path,
	}
}

// GetToken implements the TokenLocation interface.
func (k *KubernetesAuth) GetToken(client *api.Client) (string, error) {
	token, err := ioutil.ReadFile(k.path)
	if err != nil {
		return "", err
	}

	secret, err := client.Logical().Write("auth/kubernetes/login", map[string]interface{}{
		"jwt":  string(token),
		"role": k.role,
	})
	if err != nil {
		return "", err
	}

	return secret.Auth.ClientToken, nil
}
