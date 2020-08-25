package metrics

import (
	"hash/maphash"
	"sync"
)

var secureHash = &secureHasher{}

type secureHasher struct {
	h maphash.Hash
	m sync.Mutex
}

func (sh *secureHasher) GetHash(s []string) uint64 {
	sh.m.Lock()
	defer sh.m.Unlock()

	for _, part := range s {
		sh.h.WriteString(part)
	}

	return sh.h.Sum64()
}
