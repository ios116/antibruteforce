package entities

import (
	"time"
)

type KindBucket string

const (
	Login    KindBucket = "login"
	Password            = "password"
	Ip                  = "ip"
)

type Hash struct {
	Kind KindBucket
	Key  string
}

func NewHash(kind KindBucket, key string) *Hash {
	return &Hash{Kind: kind, Key: key}
}

type StoreManager interface {
	Add(key *Hash, bucket *Bucket) error
	Delete(hash *Hash) error
	Get(key *Hash) (*Bucket, error)
}

type Bucket struct {
	Marker   int
	Duration time.Duration
}

// NewBucket with a callback chanel. The chanel send message for delete bucket from storage
func NewBucket(marker int, duration time.Duration, hash *Hash, callback chan *Hash) *Bucket {
	time.AfterFunc(duration, func() {
		callback <- hash
	})
	return &Bucket{Marker: marker, Duration: duration}
}

// Counter subtract one from bucket marker
func (b *Bucket) Counter() bool {
	b.Marker -= 1
	if b.Marker < 0 {
		b.Marker = 0
		return false
	}
	return true
}
