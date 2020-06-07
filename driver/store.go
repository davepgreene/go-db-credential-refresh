package driver

// Store represents a mechanism for retrieving Credentials.
type Store interface {
	Get() (Credentials, error)
	Refresh() (Credentials, error)
}

// Credentials represents an abstraction over a username and password.
type Credentials interface {
	GetUsername() string
	GetPassword() string
}
