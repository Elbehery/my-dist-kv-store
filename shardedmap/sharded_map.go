package shardedmap

import (
	"crypto/sha1"
	"sync"
)

type shard struct {
	sync.RWMutex
	m map[string]interface{}
}

type ShardedMap []*shard

func NewShardedMap(nShards int) ShardedMap {
	shardedMap := make([]*shard, 0, nShards)
	i := 0
	for i < nShards {
		s := &shard{
			RWMutex: sync.RWMutex{},
			m:       map[string]interface{}{},
		}
		shardedMap = append(shardedMap, s)
		i++
	}
	return shardedMap
}

func (s ShardedMap) Get(Key string) interface{} {
	shard := s.getShard(Key)
	shard.RLock()
	defer shard.RUnlock()
	val, ok := shard.m[Key]
	if !ok {
		return nil
	}
	return val
}

func (s ShardedMap) Set(Key string, val interface{}) {
	shard := s.getShard(Key)
	shard.Lock()
	defer shard.Unlock()
	shard.m[Key] = val
}

func (s ShardedMap) Keys() []string {
	keys := make([]string, 0)
	mux := sync.Mutex{}
	wg := sync.WaitGroup{}

	wg.Add(len(s))
	for _, sh := range s {
		go func(s *shard) {
			defer wg.Done()
			s.RLock()
			defer s.RUnlock()
			for k := range s.m {
				mux.Lock()
				keys = append(keys, k)
				mux.Unlock()
			}
		}(sh)
	}
	wg.Wait()

	return keys
}

func (s ShardedMap) getShardIndex(Key string) int {
	chkSum := sha1.Sum([]byte(Key))
	hash := int(chkSum[17])
	return hash % len(s)
}

func (s ShardedMap) getShard(Key string) *shard {
	idx := s.getShardIndex(Key)
	return s[idx]
}
