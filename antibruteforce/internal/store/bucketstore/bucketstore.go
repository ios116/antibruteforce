package bucketstore

import (
	"antibruteforce/internal/domain/entities"
	"antibruteforce/internal/domain/exceptions"
	"sync"
)

// BucketStore
type BucketStore struct {
	mux     sync.Mutex
	Buckets map[*entities.Hash]*entities.Bucket
}

func NewBucketStore() *BucketStore {
	return &BucketStore{Buckets: make(map[*entities.Hash]*entities.Bucket)}
}

// Add bucket with key
func (st *BucketStore) Add(hash *entities.Hash, bucket *entities.Bucket) error {
	st.mux.Lock()
	st.Buckets[hash] = bucket
	st.mux.Unlock()
	return nil
}

// Delete bucket by key
func (st *BucketStore) Delete(hash *entities.Hash) error {
	st.mux.Lock()
	delete(st.Buckets, hash)
	st.mux.Unlock()
	return nil
}

// Get bucket by key
func (st *BucketStore) Get(hash *entities.Hash) (*entities.Bucket, error) {
	st.mux.Lock()
	defer st.mux.Unlock()
	bk, ok := st.Buckets[hash]
	if !ok {
		return nil, exceptions.NilValue
	}
	return bk, nil
}
