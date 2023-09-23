package util

import (
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

// Reference: https://stackoverflow.com/a/72585536

var entropyPool = sync.Pool{
	New: func() any {
		entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
		return entropy
	},
}

func NewID() string {
	e := entropyPool.Get().(*ulid.MonotonicEntropy)
	s := ulid.MustNew(ulid.Timestamp(time.Now()), e).String()
	entropyPool.Put(e)
	return s
}
