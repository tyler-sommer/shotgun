// Package config provides application configuration.
package config

import (
	"code.google.com/p/go.crypto/ssh"
	"github.com/tyler-sommer/shotgun/database"

	"io/ioutil"
)

// Config is a sort of factory for creating ready-to-use application components.
type Config struct {
	databaseFile string
	privateKeyFile string

	authMethods []ssh.AuthMethod
}

// New allocates a new Config basedon the given parameters.
func New(databaseFile, privateKeyFile string) (*Config, error) {
	c := &Config{databaseFile, privateKeyFile, make([]ssh.AuthMethod, 0)}
	err := c.populateAuthMethods()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) populateAuthMethods() error {
	privKeyText, err := ioutil.ReadFile(c.privateKeyFile)
	if err != nil {
		return err
	}

	privKey, err := ssh.ParseRawPrivateKey(privKeyText)
	if err != nil {
		return err
	}

	signer, err := ssh.NewSignerFromKey(privKey)
	if err != nil {
		return err
	}

	c.authMethods = []ssh.AuthMethod{
		ssh.PublicKeys(signer),
	}

	return nil
}

// NewClientConfig creates a ssh.ClientConfig with the given user and
// any auth methods defined in the Config.
func (c *Config) NewClientConfig(user string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: user,
		Auth: c.authMethods,
	}
}

// NewDatabaseManager creates a database.Manager with the configured
// database file.
func (c *Config) NewDatabaseManager() (*database.Manager, error) {
	return database.New(c.databaseFile)
}
