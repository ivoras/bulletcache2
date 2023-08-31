package types

import "github.com/ivoras/bulletcache2/util"

type BulletCacheBucket struct {
	util.WithRWMutex
	Items map[CacheKey]CacheRecord
}

func (bcb *BulletCacheBucket) GetByKey(key CacheKey) (rec CacheRecord) {
	bcb.WithRWMutex.RLock()
	rec = bcb.Items[key]
	bcb.WithRWMutex.RUnlock()
	return
}
