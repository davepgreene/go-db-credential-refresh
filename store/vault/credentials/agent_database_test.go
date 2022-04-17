package vaultcredentials

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/davepgreene/go-db-credential-refresh/store"
)

const (
	username = "foo"
	password = "bar"
)

var testMapper Mapper = func(s string) (*store.Credential, error) {
	creds := strings.Split(s, ":")
	if len(creds) == 2 {
		return &store.Credential{
			Username: creds[0],
			Password: creds[1],
		}, nil
	}

	return nil, errors.New("mapping function failed")
}

func TestNewAgentDatabaseCredentials(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmpfile.Name())

	contents := fmt.Sprintf("%s:%s", username, password)

	if _, err = tmpfile.Write([]byte(contents)); err != nil {
		t.Error(err)
	}
	if err = tmpfile.Close(); err != nil {
		t.Error(err)
	}

	adc := NewAgentDatabaseCredentials(testMapper, tmpfile.Name())

	credStr, err := adc.GetCredentials(context.Background(), nil)
	if err != nil {
		t.Error(err)
	}

	if credStr != contents {
		t.Errorf("expected credential string to equal '%s' but got '%s' instead", contents, credStr)
	}

	mappedCreds, err := adc.Map(credStr)
	if err != nil {
		t.Error(err)
	}

	if mappedCreds.GetUsername() != username {
		t.Errorf("expected username to be '%s' but got '%s' instead", username, mappedCreds.GetUsername())
	}

	if mappedCreds.GetPassword() != password {
		t.Errorf("expected password to be '%s' but got '%s' instead", password, mappedCreds.GetPassword())
	}
}

func TestNewAgentDatabaseCredentialsFailedFileRead(t *testing.T) {
	adc := NewAgentDatabaseCredentials(testMapper, "")
	credStr, err := adc.GetCredentials(context.Background(), nil)
	if err == nil {
		t.Error("expected an error but didn't get one")
	}
	if credStr != "" {
		t.Errorf("expected an empty output but got %s", credStr)
	}
}

func TestNewAgentDatabaseCredentialsFailedMapper(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmpfile.Name())

	contents := "foo:bar:baz"

	if _, err = tmpfile.Write([]byte(contents)); err != nil {
		t.Error(err)
	}
	if err = tmpfile.Close(); err != nil {
		t.Error(err)
	}

	adc := NewAgentDatabaseCredentials(testMapper, tmpfile.Name())

	credStr, err := adc.GetCredentials(context.Background(), nil)
	if err != nil {
		t.Error(err)
	}

	mappedCreds, err := adc.Map(credStr)
	if err == nil {
		t.Error("expected an error but didn't get one")
	}

	if mappedCreds != nil {
		t.Errorf("expected a nil credential but got %v instead", mappedCreds)
	}
}
