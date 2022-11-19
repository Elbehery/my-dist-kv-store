package shardedmap

import (
	"crypto/sha1"
	"sync"
)

type Shard struct {
	sync.RWMutex
	m map[string]interface{}
}

type ShardedMap []*Shard

func NewShardedMap(shards int) ShardedMap {
	smap := make([]*Shard, 0, shards)
	i := 0
	for i < shards {
		s := &Shard{
			RWMutex: sync.RWMutex{},
			m:       map[string]interface{}{},
		}
		smap = append(smap, s)
		i++
	}
	return smap
}

func (s ShardedMap) Get(key string) interface{} {
	shard := s.getShard(key)
	shard.RLock()
	defer shard.RUnlock()
	val, ok := shard.m[key]
	if !ok {
		return nil
	}
	return val
}

func (s ShardedMap) Set(key string, val interface{}) {
	shard := s.getShard(key)
	shard.Lock()
	defer shard.Unlock()
	shard.m[key] = val
}

func (s ShardedMap) Keys() []string {
	keys := []string{}
	mux := sync.Mutex{}
	wg := sync.WaitGroup{}

	wg.Add(len(s))
	for _, shard := range s {
		go func(s *Shard) {
			defer wg.Done()
			s.RLock()
			defer s.RUnlock()
			for k := range s.m {
				mux.Lock()
				keys = append(keys, k)
				mux.Unlock()
			}
		}(shard)
	}

	wg.Wait()
	return keys
}
func (s ShardedMap) getShardIndex(key string) int {
	chksum := sha1.Sum([]byte(key))
	hash := int(chksum[17])
	return hash % len(s)
}

func (s ShardedMap) getShard(key string) *Shard {
	idx := s.getShardIndex(key)
	return s[idx]
}
