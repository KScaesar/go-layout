package utility

import (
	"sync"
	"time"
)

type Singleflight struct {
	store sync.Map
}

// Do execute and returns the results of the given function,
// making sure that only one execution is in-flight for a given key at a time.
//
// If a duplicate comes in, the duplicate caller waits for the
// original to complete and receives the same results.
// The return value shared indicates whether val was given to multiple callers.
//
// https://pkg.go.dev/golang.org/x/sync/singleflight#Group.Do
func (group *Singleflight) Do(key string, fn func() (val any, err error)) (val any, err error, shared bool) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	wait := func() (any, error) {
		wg.Wait()
		return val, err
	}

	Wait, loaded := group.store.LoadOrStore(key, wait)
	if loaded {
		val, err = Wait.(func() (any, error))()
		return val, err, true
	}

	defer wg.Done()

	val, err = fn()
	return val, err, false
}

// Forget tells the singleflight to forget about a key.
// Future calls to Do for this key will call the function
// rather than waiting for an earlier call to complete.
func (group *Singleflight) Forget(key string) {
	group.store.Delete(key)
}

// Expire schedules a Forget operation for a given key after a specified duration.
//
// Note:
// Singleflight is designed to ensure that only one execution is in flight for a given key at a time.
// It should not be used as a local cache.
func (group *Singleflight) Expire(key string, t time.Duration) {
	if t == 0 {
		group.Forget(key)
	} else {
		time.AfterFunc(t, func() {
			group.Forget(key)
		})
	}
}
