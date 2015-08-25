package data

import (
	"github.com/boltdb/bolt"

	"os"
)

type DatabaseManager struct {
	db *bolt.DB
	tx *bolt.Tx

	buckets map[string]*bolt.Bucket
}

func New(dbFile string) (*DatabaseManager, error) {
	db, err := bolt.Open(dbFile, os.FileMode(0600), nil)
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}

	return &DatabaseManager{db, tx, make(map[string]*bolt.Bucket)}, nil
}

func (dbm *DatabaseManager) bucket(name string) (*bolt.Bucket, error) {
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

func (dbm *DatabaseManager) reset() error {
	dbm.tx = nil
	dbm.buckets = make(map[string]*bolt.Bucket)

	tx, err := dbm.db.Begin(true)
	if err != nil {
		return err
	}

	dbm.tx = tx

	return nil
}

func (dbm *DatabaseManager) Commit() error {
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

func (dbm *DatabaseManager) Rollback() error {
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

func (dbm *DatabaseManager) NewServerRepository() (*serverRepository, error) {
	return &serverRepository{dbm}, nil
}
