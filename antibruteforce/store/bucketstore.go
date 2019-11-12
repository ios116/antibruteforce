package store

import (
	"antibruteforce/domain/entities"
	"antibruteforce/domain/exceptions"
	"sync"
)

// BucketStore
type BucketStore struct {
	sync.Mutex
	Elements map[string]*entities.Bucket
}

func NewBucketStore() *BucketStore {
	return &BucketStore{Elements: make(map[string]*entities.Bucket)}
}

// Add bucket with key
func (st *BucketStore) Add(key string, bucket *entities.Bucket)  error {
	st.Lock()
	st.Elements[key] = bucket
	st.Unlock()
	return nil
}

// Delete bucket by key
func (st *BucketStore) Delete(key string) error {
	st.Lock()
	delete(st.Elements, key)
	st.Unlock()
	return nil
}

// Get bucket by key
func (st *BucketStore) Get(key string) (*entities.Bucket,error) {
	st.Lock()
	bk, ok := st.Elements[key]
	if !ok {
		return nil, exceptions.BucketsNil
	}
	st.Unlock()
	return bk, nil
}
