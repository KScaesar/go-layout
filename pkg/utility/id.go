package utility

import (
	cryptoRand "crypto/rand"

	"github.com/oklog/ulid/v2"
)

var entropyPool = NewPool(func() *ulid.MonotonicEntropy {
	return ulid.Monotonic(cryptoRand.Reader, 0)
})

func NewUlid() string {
	entropy := entropyPool.Get()
	id := ulid.MustNew(ulid.Now(), entropy)
	entropyPool.Put(entropy)
	return id.String()
}
