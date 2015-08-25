package config

import (
	"code.google.com/p/go.crypto/ssh"
	"github.com/tyler-sommer/shotgun/database"

	"io/ioutil"
)

type Config struct {
	databaseFile string
	privateKeyFile string

	authMethods []ssh.AuthMethod
}

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

func (c *Config) NewClientConfig(user string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: user,
		Auth: c.authMethods,
	}
}

func (c *Config) NewDatabaseManager() (*database.Manager, error) {
	return database.New(c.databaseFile)
}
