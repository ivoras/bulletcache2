package types

import "github.com/ivoras/utime"

const MaxTagsPerRecord = 8

type CacheRecord struct {
	MTime      utime.Time
	Flags      uint16
	Key        CacheKey
	TagsLength uint8
	Tags       [MaxTagsPerRecord]CacheTag
	Data       []byte
}

// HasTags returns true if the record is tagged with all of the given tags
func (r *CacheRecord) HasTags(tags ...CacheTag) (allPresent bool) {
	for i := 0; i < int(r.TagsLength); i++ {
		recordTag := r.Tags[i]
		found := false
		for _, queryTag := range tags {
			if queryTag == recordTag {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// HasTag returns true if the record is tagged with the given tag
func (r *CacheRecord) HasTag(tag CacheTag) (present bool) {
	for i := 0; i < int(r.TagsLength); i++ {
		if r.Tags[i] == tag {
			return true
		}
	}
	return false
}
