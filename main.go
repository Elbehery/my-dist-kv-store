package main

import (
	"fmt"
	"playground/my-dist-kv-store/shardedmap"
)

func main() {

	shardedMap := shardedmap.NewShardedMap(5)
	shardedMap.Set("alpha", 1)
	shardedMap.Set("beta", 2)
	shardedMap.Set("gamma", 3)
	fmt.Println(shardedMap.Get("alpha"))
	fmt.Println(shardedMap.Get("beta"))
	fmt.Println(shardedMap.Get("gamma"))
	keys := shardedMap.Keys()
	for _, k := range keys {
		fmt.Println(k)
	}
}
