// Package database provides an abstraction layer between the application and
// the underlying storage.
package database

import (
	"github.com/boltdb/bolt"

	"os"
)

// Manager maintains transactions and handles repository creation.
type Manager struct {
	db *bolt.DB
	tx *bolt.Tx

	buckets map[string]*bolt.Bucket
}

// New creates a Manager using the given file.
func New(dbFile string) (*Manager, error) {
	db, err := bolt.Open(dbFile, os.FileMode(0600), nil)
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}

	return &Manager{db, tx, make(map[string]*bolt.Bucket)}, nil
}

// bucket returns a Bucket with the given name for the current active transaction.
func (dbm *Manager) bucket(name string) (*bolt.Bucket, error) {
	if val, ok := dbm.buckets[name]; ok {
		return val, nil
	}

	b, err := dbm.tx.CreateBucketIfNotExists([]byte("Server"))
	if err != nil {
		return nil, err
	}

	dbm.buckets[name] = b

	return b, nil
}

// reset clears the Manager and opens a new transaction.
func (dbm *Manager) reset() error {
	dbm.tx = nil
	dbm.buckets = make(map[string]*bolt.Bucket)

	tx, err := dbm.db.Begin(true)
	if err != nil {
		return err
	}

	dbm.tx = tx

	return nil
}

// Commit attempts to persist any changes to the underlying database.
func (dbm *Manager) Commit() error {
	err := dbm.tx.Commit()
	if err != nil {
		return err
	}

	err = dbm.reset()
	if err != nil {
		return err
	}

	return nil
}

// Rollback resets any changes made in the current transaction and starts a new one.
func (dbm *Manager) Rollback() error {
	err := dbm.tx.Rollback()
	if err != nil {
		return err
	}

	err = dbm.reset()
	if err != nil {
		return err
	}

	return nil
}

// NewServerRepository allocates a fully-wired ServerRepository.
func (dbm *Manager) NewServerRepository() (*ServerRepository, error) {
	return &ServerRepository{dbm}, nil
}
