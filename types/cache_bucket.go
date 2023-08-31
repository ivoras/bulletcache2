package types

import (
	"github.com/ivoras/bulletcache2/errors"
	"github.com/ivoras/bulletcache2/util"
	"github.com/ivoras/utime"
)

type BulletCacheBucket struct {
	util.WithRWMutex
	Items map[CacheKey]CacheRecord
}

// GetByKey returns a record identified by the key
func (bcb *BulletCacheBucket) GetByKey(key CacheKey) (rec CacheRecord) {
	bcb.WithRWMutex.RLock()
	rec = bcb.Items[key]
	bcb.WithRWMutex.RUnlock()
	return
}

// Set sets a record in the cache bucket
func (bcb *BulletCacheBucket) Set(key CacheKey, value []byte, onlyIfMissing bool, flags uint32, expTime uint32) (err error) {
	key.Sanitize()

	bcb.WithRWMutex.Lock()
	defer bcb.WithRWMutex.Unlock()

	_, found := bcb.Items[key]
	if found && onlyIfMissing {
		return errors.ErrKeyAlreadyExists
	}

	var recordExpTime utime.Time
	if expTime == 0 {
		// Special value, leaves recordExpTime as 0
	} else if expTime < 1_000_000_000 {
		recordExpTime = utime.Now() + utime.Time(expTime)
	} else {
		recordExpTime = utime.Time(expTime)
	}

	bcb.Items[key] = CacheRecord{
		ExpTime: recordExpTime,
		Key:     key,
		Value:   value,
	}

	return nil
}
