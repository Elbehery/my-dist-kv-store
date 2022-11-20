package kvstore

import (
	"errors"
	"sync"
)

var (
	ErrKeyNotExist = errors.New("key does not exist")
)

type KVStore struct {
	sync.RWMutex
	s map[string]string
}

func NewKVStore() *KVStore {
	s := KVStore{s: map[string]string{}}
	return &s
}

func (s *KVStore) Put(key, value string) error {
	s.Lock()
	defer s.Unlock()
	s.s[key] = value
	return nil
}

func (s *KVStore) Get(key string) (string, error) {
	s.RLock()
	defer s.RUnlock()
	val, ok := s.s[key]
	if !ok {
		return "", ErrKeyNotExist
	}
	return val, nil
}

func (s *KVStore) Delete(key string) error {
	s.Lock()
	defer s.Unlock()
	delete(s.s, key)
	return nil
}
