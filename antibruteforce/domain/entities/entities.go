package entities

import (
	"time"
)

type Kind string

const (
	Login    Kind = "login"
	Password      = "password"
	Ip            = "ip"
)

type Hash struct {
	Kind Kind
	Key string
}

func NewHash(kind Kind, key string) *Hash {
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

//func (b *Bucket) CreateKey(key string, kind Kind) string {
//	var buffer bytes.Buffer
//	buffer.WriteString(key)
//	buffer.WriteRune(':')
//	buffer.WriteString(string(kind))
//	return buffer.String()
//}

// Counter subtract one from bucket marker
func (b *Bucket) Counter() bool {
	b.Marker -= 1
	if b.Marker < 0 {
		b.Marker = 0
		return false
	}
	return true
}
