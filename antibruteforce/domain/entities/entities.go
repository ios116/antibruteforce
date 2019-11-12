package entities

import (
	"time"
)

type Kind int

const (
	Login Kind = iota
	Password
	Ip
)

type StoreManager interface {
	Add(key string, bucket *Bucket) error
	Delete(key string) error
	Get(key string) (*Bucket, error)
}

type Bucket struct {
	Marker   int
	Duration time.Duration
}

// NewBucket with a callback chanel. The chanel send message for delete bucket from storage
func NewBucket(marker int, duration time.Duration, key string, delete chan string) *Bucket {
	time.AfterFunc(duration, func() {
		delete <- key
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
