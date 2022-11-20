package kvstore

import (
	"errors"
	"sync"
)

var (
	ErrKeyNotExist = errors.New("Key does not exist")
)

type KVStore struct {
	sync.RWMutex
	s map[string]string
}

func NewKVStore() *KVStore {
	s := KVStore{s: map[string]string{}}
	return &s
}

func (s *KVStore) Put(Key, value string) error {
	s.Lock()
	defer s.Unlock()
	s.s[Key] = value
	return nil
}

func (s *KVStore) Get(Key string) (string, error) {
	s.RLock()
	defer s.RUnlock()
	val, ok := s.s[Key]
	if !ok {
		return "", ErrKeyNotExist
	}
	return val, nil
}

func (s *KVStore) Delete(Key string) error {
	s.Lock()
	defer s.Unlock()
	delete(s.s, Key)
	return nil
}
