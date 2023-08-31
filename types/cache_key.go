package types

import (
	"bytes"

	"github.com/zeebo/xxh3"
)

const MaxKeyLength = 255

type CacheKey struct {
	KeyLength uint8
	Key       [MaxKeyLength]byte
}

func (ck CacheKey) Equals(ck2 CacheKey) bool {
	return ck.KeyLength == ck2.KeyLength && bytes.Equal(ck.Key[0:ck.KeyLength], ck2.Key[0:ck.KeyLength])
}

func (ck *CacheKey) Sanitize() {
	for i := ck.KeyLength; i < MaxKeyLength; i++ {
		ck.Key[i] = 0
	}
}

func (ck *CacheKey) Hash1() (result uint8) {
	return uint8(xxh3.Hash(ck.Key[0:ck.KeyLength]))
}
