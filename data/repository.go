package data

import (
	"github.com/tyler-sommer/shotgun/model"

	"math/rand"
	"time"
	"encoding/json"
	"errors"
)

var KeyNotFoundError = errors.New("Unable to locate record with the given key")

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func genKey() string {
	n := 32
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

type serverRepository struct {
	dbm *DatabaseManager
}

func transformServer(data []byte) (model.Server, error) {
	s := model.Server{}
	err := json.Unmarshal(data, &s)

	return s, err
}

func reverseTransformServer(s model.Server) ([]byte, error) {
	res, err := json.Marshal(&s)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (repo *serverRepository) All() ([]model.Server, error) {
	bucket, err := repo.dbm.bucket("Server")
	if err != nil {
		return nil, err
	}

	servers := make([]model.Server, 0)

	err = bucket.ForEach(func(k, v []byte) error {
		s, err := transformServer(v)
		if err != nil {
			return err
		}

		s.SetKey(string(k))

		servers = append(servers, s)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (repo *serverRepository) Find(key string) (model.Server, error) {
	bucket, err := repo.dbm.bucket("Server")
	if err != nil {
		return model.Server{}, err
	}

	val := bucket.Get([]byte(key))
	if len(val) == 0 {
		return model.Server{}, KeyNotFoundError
	}

	s, err := transformServer(val)
	s.SetKey(key)

	return s, err
}

func (repo *serverRepository) Save(s model.Server) error {
	bucket, err := repo.dbm.bucket("Server")
	if err != nil {
		return err
	}

	key := s.Key()
	if len(key) == 0 {
		key = genKey()
		s.SetKey(key)
	}

	val, err := reverseTransformServer(s)
	if err != nil {
		return err
	}

	return bucket.Put([]byte(key), val)
}

func (repo *serverRepository) Delete(key string) error {
	bucket, err := repo.dbm.bucket("Server")
	if err != nil {
		return err
	}

	err = bucket.Delete([]byte(key))
	if err != nil {
		return err
	}

	return nil
}
