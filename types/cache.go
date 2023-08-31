package types

type BulletCache struct {
	// We take the first byte of the uint32 hash of the key to find out which
	// lockable bucket contains the key. This is an unlocked operation because
	// this particular bucket list is static.
	Buckets [256]BulletCacheBucket
}

func NewBulletCache() *BulletCache {
	return &BulletCache{}
}

func (bc *BulletCache) GetByKey(key CacheKey) (rec CacheRecord) {
	return bc.Buckets[key.Hash1()].GetByKey(key)
}
