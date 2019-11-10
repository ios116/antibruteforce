package store

import (
	"antibruteforce/domain"
	"sync"
)



type BucketStore struct {
	sync.Mutex
	Login map[string]*domain.Bucket
}

func (st *BucketStore) Add(key string, bucket *domain.Bucket) {
	st.Lock()
	st.Login[key] = bucket
	st.Unlock()
}

func (st *BucketStore) Delete(key string) {
	st.Lock()
	delete(st.Login, key)
	st.Unlock()
}
func (st *BucketStore) Get(key string) *domain.Bucket {
	st.Lock()
	bk, ok := st.Login[key]
	if !ok {
		return nil
	}
	st.Unlock()
	return bk
}
