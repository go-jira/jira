package keyring

import (
	"errors"
	"sync"
)

var (
	// ErrNotFound means the requested password was not found
	ErrNotFound = errors.New("keyring: Password not found")
	// ErrNoDefault means that no default keyring provider has been found
	ErrNoDefault = errors.New("keyring: No suitable keyring provider found (check your build flags)")

	providerInitOnce  sync.Once
	defaultProvider   provider
	providerInitError error
)

// provider provides a simple interface to keychain sevice
type provider interface {
	Get(service, username string) (string, error)
	Set(service, username, password string) error
}

func setupProvider() (provider, error) {
	providerInitOnce.Do(func() {
		defaultProvider, providerInitError = initializeProvider()
	})

	if providerInitError != nil {
		return nil, providerInitError
	} else if defaultProvider == nil {
		return nil, ErrNoDefault
	}
	return defaultProvider, nil
}

// Get gets the password for a paricular Service and Username using the
// default keyring provider.
func Get(service, username string) (string, error) {
	p, err := setupProvider()
	if err != nil {
		return "", err
	}

	return p.Get(service, username)
}

// Set sets the password for a particular Service and Username using the
// default keyring provider.
func Set(service, username, password string) error {
	p, err := setupProvider()
	if err != nil {
		return err
	}

	return p.Set(service, username, password)
}
