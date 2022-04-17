package vaultauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	authv1 "k8s.io/api/authentication/v1"

	"github.com/davepgreene/go-db-credential-refresh/store/vault/vaulttest"
)

var (
	testCACert = `-----BEGIN CERTIFICATE-----
MIIC5zCCAc+gAwIBAgIBATANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDEwptaW5p
a3ViZUNBMB4XDTE5MDEwNTE4MDkxNFoXDTI5MDEwMzE4MDkxNFowFTETMBEGA1UE
AxMKbWluaWt1YmVDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMNw
BYLyqYEVm4vbik0NibQ7+G414dPUUZc3UCScrvMYASa+Krcc8J8Ic0TeDdsluiYs
hujALbu+LtFNYeIpMBgZPUaBVOtSrnBe9ieG0XZmxDa303uz2awzYkivWab58Tsx
RLojX2z4ZJUXhb1m6VN96x07tf4MujnQgmfm3GZ/cMn/BUaTSJOKXiKTDTys6dbz
U3UyvQnxP9QkWloU6HICqPObzpY/kkLdsOWPfiGn2lINZ/9zkeW8Qe9QalKRuGI2
2+ZWOTZyREvfln/3LML8q9kAmk54NMtSG3mGCgDOL+HsRVNnqZC2QBHmHJGD+2nz
z6C1iSV0W4ZgDR0HSIMCAwEAAaNCMEAwDgYDVR0PAQH/BAQDAgKkMB0GA1UdJQQW
MBQGCCsGAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3
DQEBCwUAA4IBAQC1epomLINu4JVBCZpaHDpWg5cAHyxAVggKihCP+hIsOkfsmyTa
RPCLNASmW/PbDDzJKAzQaC23KDFW9WCqr1MlgsJhZMW8tkiPegL18DGxupjwzIIM
meiZoPBEpFGz0JVhfu0FMIVbvKjhuBgbZd3rKZEFHZMer7L+ZZ2Pd/5UY1s7oslq
Z938fecvWwIQLHE+Jar1KqvdlLlP798+w7G5de3gIN0svflcpbd8+w3X7h+dzouu
qevae2NSZJ5r8Fo5Ch3sI63c6GCoUaMM5Ho7mHUM32BeGxy99Z3G6364akR3I819
qQYZl8EZf4Jznaes/XFP0Yb+IhGXBoR9Ib+I
-----END CERTIFICATE-----`

	jwtData     = `eyJhbGciOiJSUzI1NiIsImtpZCI6IlpJOEY4RHVoMktrY0JxTjhGSGxyMEhER2l2OEtFR2xFSnlITUZRc1UwZ28ifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6InZhdWx0LWF1dGgtdG9rZW4tdmQ0bjQiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoidmF1bHQtYXV0aCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjlhZjM3NjRlLWZmZDMtNDJiZC1hZjVkLTE2MzUwZTM0NjkyYyIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OnZhdWx0LWF1dGgifQ.ZfkKFqeAIaNXmk-i7LwrXXoOIjv4WlQ1gHFOXHpSo0Wdq16KKu1VOnCkzUh9bIApL5pIXZu4-eYwP2SwokRafXBY_5znqvXoI3F1fxmw25jBT9ZeyDEKZOxyO7mtHnh7LZQ_pBUPPflClhAwacbBrTjnIpHoiWq-Z1_BeuenlRdBYQYjdXEOPK-W1bFbCqx4hq_x91v-JMAcJqQUf0ZSY3jU-vcAOmFfv_0S4K2_syUyfkYVPr_pX-0wOvwkv0nDhV-fhqux51onQyYDd_gejvjGvviDJcbXxT4sIYgbS8IKtRwI3lAhpQQyuaQbVI6DKASs9z-jvvg0VO7T2FMFIw`
	jwtUID      = "9af3764e-ffd3-42bd-af5d-16350e34692c"
	jwtUsername = "system:serviceaccount:default:vault-auth"
	jwtGroups   = []string{
		"system:serviceaccounts",
		"system:serviceaccounts:default",
		"system:authenticated",
	}
)

// this HandlerFunc mocks out an response from k8s's tokenreviews endpoint.
func tokenReviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/apis/authentication.k8s.io/v1/tokenreviews" {
		w.WriteHeader(404)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer r.Body.Close()

	var tr authv1.TokenReview
	if err := json.Unmarshal(body, &tr); err != nil {
		w.WriteHeader(500)
		return
	}

	tr.Status = authv1.TokenReviewStatus{
		Authenticated: true,
		User: authv1.UserInfo{
			UID:      jwtUID,
			Username: jwtUsername,
			Groups:   jwtGroups,
		},
	}

	json.NewEncoder(w).Encode(tr) //nolint:errcheck
}

func TestKubernetesAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(tokenReviewHandler))
	defer srv.Close()

	ln, client := vaulttest.CreateTestVault(t, nil)
	defer ln.Close()

	ctx := context.Background()

	if _, err := client.Logical().WriteWithContext(ctx, "sys/auth/kubernetes", map[string]interface{}{
		"type": "kubernetes",
	}); err != nil {
		t.Error(err)
	}

	if _, err := client.Logical().WriteWithContext(ctx, "auth/kubernetes/config", map[string]interface{}{
		"kubernetes_host":    srv.URL,
		"kubernetes_ca_cert": testCACert,
	}); err != nil {
		t.Error(err)
	}

	role := "example"
	userName := strings.Split(jwtUsername, ":")

	if _, err := client.Logical().WriteWithContext(ctx, fmt.Sprintf("auth/kubernetes/role/%s", role), map[string]interface{}{
		"bound_service_account_names":      userName[len(userName)-1],
		"bound_service_account_namespaces": "default",
	}); err != nil {
		t.Error(err)
	}

	tmpfile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err = tmpfile.Write([]byte(jwtData)); err != nil {
		t.Error(err)
	}
	if err = tmpfile.Close(); err != nil {
		t.Error(err)
	}

	a := NewKubernetesAuth(role, tmpfile.Name())
	token, err := a.GetToken(ctx, client)
	if err != nil {
		t.Error(err)
	}

	if token == "" {
		t.Error("expected a token but didn't get one")
	}

	// Verify the token
	client.SetToken(token)
	resp, err := client.Logical().ReadWithContext(ctx, tokenSelfLookupPath)
	if err != nil {
		t.Error(err)
	}

	if resp == nil {
		t.Fatal("expected a valid response")
	}

	if resp.Data == nil {
		t.Fatal("expected response to have data")
	}

	path, ok := resp.Data["path"]
	if !ok {
		t.Error("expected 'path' to be in auth response data")
	}

	if path != "auth/kubernetes/login" {
		t.Errorf("expected 'path' to be k8s login path but got %s instead", path)
	}
}

func TestKubernetesAuthFileError(t *testing.T) {
	p := "/foo/bar/baz"
	k := NewKubernetesAuth("role", p)
	client, err := api.NewClient(nil)
	if err != nil {
		t.Error(err)
	}
	_, err = k.GetToken(context.Background(), client)
	if err == nil {
		t.Error("expected an error but didn't get one")
	}

	pathErr := &os.PathError{
		Op:   "open",
		Path: p,
		Err:  errors.New("no such file or directory"),
	}
	if err.Error() != pathErr.Error() {
		t.Errorf("expected error to be '%v' but got '%v' instead", pathErr, err)
	}
}

func TestKubernetesAuthVaultError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc((func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})))
	defer srv.Close()

	client, err := api.NewClient(nil)
	if err != nil {
		t.Error(err)
	}

	if err = client.SetAddress(srv.URL); err != nil {
		t.Error(err)
	}

	tmpfile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err = tmpfile.Write([]byte(jwtData)); err != nil {
		t.Error(err)
	}
	if err = tmpfile.Close(); err != nil {
		t.Error(err)
	}

	k := NewKubernetesAuth("role", tmpfile.Name())
	_, err = k.GetToken(context.Background(), client)
	if err == nil {
		t.Error("expected an error but didn't get one")
	}

	respErr := &api.ResponseError{}
	if errors.As(err, &respErr) {
		if respErr.StatusCode != http.StatusNotFound {
			t.Errorf("expected to get a %d but got a %d instead", http.StatusNotFound, respErr.StatusCode)
		}

		loginURL := fmt.Sprintf("%s/v1/auth/kubernetes/login", srv.URL)
		if respErr.URL != loginURL {
			t.Errorf("expected URL to be %s but got %s instead", loginURL, respErr.URL)
		}
		if respErr.HTTPMethod != http.MethodPut {
			t.Errorf("expected method %s but got %s instead", http.MethodPut, respErr.HTTPMethod)
		}
	}
}
