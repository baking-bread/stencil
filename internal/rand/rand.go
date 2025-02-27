package rand

import (
	"math/rand"
	"sync"
	"time"
)

var (
	rnd  *rand.Rand
	once sync.Once
)

func init() {
	once.Do(func() {
		rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	})
}

func Random(limit int) int {
	if limit <= 0 {
		return 0
	}
	return rnd.Intn(limit)
}
