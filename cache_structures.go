package main

import (
	"github.com/ivoras/utime"
)

type CacheRecord struct {
	MTime     utime.Time
	Flags     uint16
	NTagsSet  uint8
	KeyLength uint8
	Tags      [4]uint32
	Key       [255]byte
}
