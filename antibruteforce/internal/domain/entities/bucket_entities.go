package entities

import (
	"time"
)

type KindBucket string

const (
	Login    KindBucket = "login"
	Password            = "password"
	IP                  = "ip"
)

// Hash is used as а key for the buckets. Includes key and type (login or password or ip)
type Hash struct {
	Kind KindBucket
	Key  string
}

// NewHash created instance of key
func NewHash(kind KindBucket, key string) *Hash {
	return &Hash{Kind: kind, Key: key}
}

// BucketStoreManager bucket store interface
type BucketStoreManager interface {
	Add(key *Hash, bucket *Bucket) error
	Delete(hash *Hash) error
	Get(key *Hash) (*Bucket, error)
}

// Bucket - time interval and restriction for this interval
type Bucket struct {
	Marker   int
	Duration time.Duration
}

// NewBucket bucket instance with a callback chanel. The chanel send message for delete bucket from storage
func NewBucket(marker int, duration time.Duration, hash *Hash, callback chan *Hash) *Bucket {
	time.AfterFunc(duration, func() {
		callback <- hash
	})
	return &Bucket{Marker: marker, Duration: duration}
}

// Counter subtract one from bucket marker
func (b *Bucket) Counter() bool {
	b.Marker--
	if b.Marker < 0 {
		b.Marker = 0
		return false
	}
	return true
}
