package awsrds

import (
	"context"
	"errors"
	"net/url"
	"reflect"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/mitchellh/mapstructure"
)

func TestStoreValidation(t *testing.T) {
	if _, err := NewStore(nil); err == nil {
		t.Error("expected an error but didn't get one")
	} else if err != errMissingConfig {
		t.Errorf("expected '%T' but got '%T' instead", errMissingConfig, err)
	}

	testCases := []struct {
		description string
		fields      map[string]interface{}
		expectedErr error
	}{
		{
			description: "missing endpoint",
			fields: map[string]interface{}{
				"region": "us-east-1",
				"user":   "bar",
			},
			expectedErr: &errMissingConfigItem{item: "endpoint"},
		},
		{
			description: "missing region",
			fields: map[string]interface{}{
				"endpoint": "foo",
				"user":     "bar",
			},
			expectedErr: &errMissingConfigItem{item: "region"},
		},
		{
			description: "missing user",
			fields: map[string]interface{}{
				"endpoint": "foo",
				"region":   "us-east-1",
			},
			expectedErr: &errMissingConfigItem{item: "user"},
		},
		{
			description: "malformed endpoint - no port",
			fields: map[string]interface{}{
				"endpoint": "foo",
				"region":   "us-east-1",
				"user":     "bar",
			},
			expectedErr: errMalformedEndpoint,
		},
		{
			description: "malformed endpoint - non-numeric port",
			fields: map[string]interface{}{
				"endpoint": "foo:bar",
				"region":   "us-east-1",
				"user":     "bar",
			},
			expectedErr: &url.Error{
				Op:  "parse",
				URL: "http://foo:bar",
				Err: errors.New(`invalid port ":bar" after host`),
			},
		},
		{
			description: "malformed endpoint - missing hostname",
			fields: map[string]interface{}{
				"endpoint": "http://:5432",
				"region":   "us-east-1",
				"user":     "bar",
			},
			expectedErr: errMalformedEndpoint,
		},
		{
			description: "missing credentials",
			fields: map[string]interface{}{
				"endpoint": "localhost:5432",
				"region":   "us-east-1",
				"user":     "bar",
			},
			expectedErr: errMissingCredentials,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			var conf Config
			err := mapstructure.Decode(testCase.fields, &conf)
			if err != nil {
				t.Fatal(err)
			}

			_, err = NewStore(&conf)
			if err == nil {
				t.Error("expected an error but didn't get one")
				return
			}

			// If we have a pointer to an error we need to compare error strings
			if reflect.ValueOf(testCase.expectedErr).Kind() == reflect.Ptr && err.Error() != testCase.expectedErr.Error() {
				t.Errorf("expected '%v' but got '%v' instead", testCase.expectedErr, err)
				return
			}

			if reflect.ValueOf(testCase.expectedErr).Kind() != reflect.Ptr && err != testCase.expectedErr {
				t.Errorf("expected '%T' but got '%T' instead", testCase.expectedErr, err)
			}
		})
	}

	if _, err := NewStore(&Config{
		Endpoint:    "http://localhost:5432",
		Region:      "us-east-1",
		User:        "dbuser",
		Credentials: aws.AnonymousCredentials{},
	}); err != nil {
		t.Errorf("expected no error but got %v instead", err)
	}
}

func TestValidStoreCanGenerateToken(t *testing.T) {
	s, err := NewStore(&Config{
		Endpoint:    "rdsmysql.cdgmuqiadpid.us-east-1.rds.amazonaws.com:5432",
		Region:      "us-east-1",
		User:        "dbuser",
		Credentials: credentials.NewStaticCredentialsProvider("foo", "bar", "baz"),
	})
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()

	creds, err := s.Get(ctx)
	if err != nil {
		t.Error(err)
	}

	if creds.GetUsername() == "" {
		t.Error("got empty username")
	}

	if creds.GetPassword() == "" {
		t.Error("got empty password")
	}
}

func TestStoreErrorsOnUnsignableCredentials(t *testing.T) {
	s, err := NewStore(&Config{
		Endpoint:    "rdsmysql.cdgmuqiadpid.us-east-1.rds.amazonaws.com:5432",
		Region:      "us-east-1",
		User:        "dbuser",
		Credentials: aws.AnonymousCredentials{},
	})
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()

	if _, err := s.Get(ctx); err == nil {
		t.Error("expected an error but didn't get one")
	}
}

func TestStoreCachesCredentials(t *testing.T) {
	s, err := NewStore(&Config{
		Endpoint:    "rdsmysql.cdgmuqiadpid.us-east-1.rds.amazonaws.com:5432",
		Region:      "us-east-1",
		User:        "dbuser",
		Credentials: credentials.NewStaticCredentialsProvider("foo", "bar", "baz"),
	})
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()

	creds, err := s.Get(ctx)
	if err != nil {
		t.Error(err)
	}

	username := creds.GetUsername()
	if username == "" {
		t.Error("got empty username")
	}

	password := creds.GetPassword()
	if password == "" {
		t.Error("got empty password")
	}

	// NOTE: This is hacky as hell but necessary because the rdsutil.BuildAuthToken has a hard-coded
	// 15 minute expiration for each signed token. To ensure we don't repeatedly generate the same signing
	// token we need to wind the clock forward past the 15 minute window.
	var patch *monkey.PatchGuard

	patch = monkey.Patch(time.Now, func() time.Time {
		patch.Unpatch()
		defer patch.Restore()
		return time.Now().Add(20 * time.Minute)
	})
	defer patch.Unpatch()

	// Second time through we should have everything cached
	cachedCreds, err := s.Get(ctx)
	if err != nil {
		t.Error(err)
	}

	if username != cachedCreds.GetUsername() {
		t.Errorf("expected username to be cached but got %s instead", cachedCreds.GetUsername())
	}
	if password != cachedCreds.GetPassword() {
		t.Errorf("expected password to be cached but got %s instead", cachedCreds.GetPassword())
	}

	// On refresh, we should have a new password
	refreshedCreds, err := s.Refresh(ctx)
	if err != nil {
		t.Error(err)
	}

	if password == refreshedCreds.GetPassword() {
		t.Error("cached password and refreshed password were the same but expected them not to be", password, refreshedCreds.GetPassword())
	}
}
