package kvstore

import "errors"

var (
	ErrKeyNotExist = errors.New("key does not exist")
)

type KVStore map[string]string

func NewKVStore() KVStore {
	s := KVStore{}
	return s
}

func (s KVStore) Put(key, value string) error {
	s[key] = value
	return nil
}

func (s KVStore) Get(key string) (string, error) {
	val, ok := s[key]
	if !ok {
		return "", ErrKeyNotExist
	}
	return val, nil
}

func (s KVStore) Delete(key string) error {
	delete(s, key)
	return nil
}
