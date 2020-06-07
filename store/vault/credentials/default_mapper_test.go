package vaultcredentials

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDefaultMapperMapsValidCredentials(t *testing.T) {
	creds, err := DefaultMapper(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password))
	if err != nil {
		t.Error(err)
	}

	if creds.GetUsername() != username {
		t.Errorf("expected username to be '%s' but got '%s' instead", username, creds.GetUsername())
	}

	if creds.GetPassword() != password {
		t.Errorf(" expected password to be '%s' but got '%s' instead", password, creds.GetPassword())
	}
}

func TestDefaultMapperGetsImproperlyFormedCredentials(t *testing.T) {
	testCases := []struct {
		description string
		input       string
		expectedErr error
	}{
		{
			description: "invalid username field name",
			input:       `{"user": "foo", "password": "bar"}`,
			expectedErr: errMissingUserName,
		},
		{
			description: "invalid password field name",
			input:       `{"username": "foo", "pass": "bar"}`,
			expectedErr: errMissingPassword,
		},
		{
			description: "invalid data structure",
			input:       `{"username": {"foo": "bar"}, "password": {"foo": "bar"}}`,
			expectedErr: errMissingUserName,
		},
		{
			description: "empty JSON object",
			input:       `{}`,
			expectedErr: errMissingUserName,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			if creds, err := DefaultMapper(testCase.input); err == nil && creds != nil {
				t.Errorf("expected an error and nil creds but got no error and %v instead", creds)
			} else if err != testCase.expectedErr {
				t.Errorf("expected a '%T' but got '%T'", testCase.expectedErr, err)
			}
		})
	}
}

func TestDefaultMapperGetsUnmarshalableString(t *testing.T) {
	for _, input := range []string{"foo bar baz", ""} {
		if creds, err := DefaultMapper(input); err == nil && creds != nil {
			t.Errorf("expected an error and nil creds but got no error and %v instead", creds)
		} else if _, ok := err.(*json.SyntaxError); !ok {
			t.Errorf("expected a 'json.SyntaxError' but got '%T' instead", err)
		}
	}
}
