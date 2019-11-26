package bucketstore

import (
	"antibruteforce/internal/domain/entities"
	"sync"
)

// BucketStore stores for bucket in memory
type BucketStore struct {
	mux     sync.Mutex
	Buckets map[entities.Hash]*entities.Bucket
}

// NewBucketStore create bucket storage storage
func NewBucketStore() *BucketStore {
	return &BucketStore{Buckets: make(map[entities.Hash]*entities.Bucket)}
}

// Add bucket with key
func (st *BucketStore) Add(hash entities.Hash, bucket *entities.Bucket) error {
	st.mux.Lock()
	st.Buckets[hash] = bucket
	st.mux.Unlock()
	return nil
}

// Delete bucket by hash
func (st *BucketStore) Delete(hash entities.Hash) error {
	st.mux.Lock()
	delete(st.Buckets, hash)
	st.mux.Unlock()
	return nil
}

// Get bucket by key
func (st *BucketStore) Get(hash entities.Hash) (*entities.Bucket, error) {
	st.mux.Lock()
	defer st.mux.Unlock()
	bk, ok := st.Buckets[hash]
	if !ok {
		return nil, nil
	}
	return bk, nil
}

// TotalBuckets total buckets amount
func (st *BucketStore) TotalBuckets() int {
	st.mux.Lock()
	total := len(st.Buckets)
	st.mux.Unlock()
	return total
}
