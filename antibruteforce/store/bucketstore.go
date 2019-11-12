package store

import (
	"antibruteforce/domain"
	"sync"
)

// BucketStore
type BucketStore struct {
	sync.Mutex
	Elements map[string]*domain.Bucket
}

func NewBucketStore() *BucketStore {
	return &BucketStore{Elements: make(map[string]*domain.Bucket)}
}

// Add bucket with key
func (st *BucketStore) Add(key string, bucket *domain.Bucket) {
	st.Lock()
	st.Elements[key] = bucket
	st.Unlock()
}

// Delete bucket by key
func (st *BucketStore) Delete(key string) {
	st.Lock()
	delete(st.Elements, key)
	st.Unlock()
}

// Get bucket by key
func (st *BucketStore) Get(key string) *domain.Bucket {
	st.Lock()
	bk, ok := st.Elements[key]
	if !ok {
		return nil
	}
	st.Unlock()
	return bk
}
