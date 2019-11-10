package domain

import (
	"fmt"
	"log"
	"time"
)

type Kind int

const (
	Login Kind = iota
	Password
	Ip
)

type StoreManager interface {
	Add(key string, bucket *Bucket)
	Delete(key string)
	Get(key string) *Bucket
}

type Bucket struct {
	Marker   int
	Duration time.Duration
}

func NewBucket(marker int, duration time.Duration, key string, delete chan string) *Bucket {
	time.AfterFunc(duration, func() {
		fmt.Println("timeout")
		delete <- key
	})
	return &Bucket{Marker: marker, Duration: duration}
}

func (b *Bucket) Counter() bool {
	b.Marker -= 1
	log.Println("marker remainder", b.Marker)
	if b.Marker < 0 {
		log.Println("marker is finished")
		return false
	}
	return true
}
