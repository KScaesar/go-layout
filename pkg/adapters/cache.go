package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/allegro/bigcache/v3"

	"github.com/KScaesar/go-layout/pkg"
)

func NewLocalCache() (*bigcache.BigCache, error) {
	config := bigcache.DefaultConfig(1 * time.Minute)
	config.Logger = log.Default()

	cache, err := bigcache.New(context.Background(), config)
	if err != nil {
		return nil, err
	}

	id := fmt.Sprintf("LocalCache(%p)", cache)
	pkg.Shutdown().AddPriorityShutdownAction(2, id, cache.Close)
	return cache, nil
}

func GetLocalCacheByType[T any](store *bigcache.BigCache, key string) (val T, Err error) {
	bData, err := store.Get(key)
	if err != nil {
		Err = fmt.Errorf("key=%q: bigcache.Get: %w: %w", key, pkg.ErrSystem, err)
		return
	}
	err = json.Unmarshal(bData, &val)
	if err != nil {
		Err = fmt.Errorf("key=%q: json.Unmarshal: %w: %w", key, pkg.ErrSystem, err)
		return
	}
	return
}

func SetLocalCacheByType[T any](store *bigcache.BigCache, key string, val *T) (Err error) {
	bData, err := json.Marshal(val)
	if err != nil {
		Err = fmt.Errorf("key=%q: json.Marshal: %w: %w", key, pkg.ErrSystem, err)
		return
	}
	err = store.Set(key, bData)
	if err != nil {
		Err = fmt.Errorf("key=%q: bigcache.Set: %w: %w", key, pkg.ErrSystem, err)
		return
	}
	return
}

// CachePolicyByQty
// This is used in scenarios involving count accumulation to decide when to use local cache.
//
// When current value is too far from target value, local cache is used to reduce database load.
// If current value is close to target value, local cache is bypassed.
//
// Return bool:
//   true  - Use local cache (current value is far from the target).
//   false - Bypass local cache (current value is close to the target).
func CachePolicyByQty(target int) func(current int) bool {
	const percent = 10
	tolerance := (target * percent) / 100

	return func(current int) bool {
		return !(target-tolerance <= current && current < target)
	}
}
