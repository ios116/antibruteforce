package store

import (
	"antibruteforce/domain/entities"
	"antibruteforce/domain/exceptions"
	"sync"
)

// BucketStore
type BucketStore struct {
	mux sync.Mutex
	Elements map[string]*entities.Bucket
}

func NewBucketStore() *BucketStore {
	return &BucketStore{Elements: make(map[string]*entities.Bucket)}
}

// Add bucket with key
func (st *BucketStore) Add(key string, bucket *entities.Bucket)  error {
	st.mux.Lock()
	st.Elements[key] = bucket
	st.mux.Unlock()
	return nil
}

// Delete bucket by key
func (st *BucketStore) Delete(key string) error {
	st.mux.Lock()
	delete(st.Elements, key)
	st.mux.Unlock()
	return nil
}

// Get bucket by key
func (st *BucketStore) Get(key string) (*entities.Bucket,error) {
	st.mux.Lock()
	defer st.mux.Unlock()
	bk, ok := st.Elements[key]
	if !ok {
		return nil, exceptions.BucketsNil
	}
	return bk, nil
}
