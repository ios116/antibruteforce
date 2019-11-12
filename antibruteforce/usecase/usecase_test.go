package usecase

import (
	"antibruteforce/config"
	"antibruteforce/domain"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockedBucketStore struct {
	mock.Mock
}

func (s *MockedBucketStore) Add(key string, bucket *domain.Bucket) {}
func (s *MockedBucketStore) Delete(key string)                     {}
func (s *MockedBucketStore) Get(key string) *domain.Bucket {
	args := s.Called(key)
	return args.Get(0).(*domain.Bucket)
}

func TestCheck(t *testing.T) {
	testObj := new(MockedBucketStore)
	bucket := &domain.Bucket{
		Marker:   0,
		Duration: 5000,
	}
	settings := config.NewSettings()
	testObj.On("Get", "admin").Return(bucket)
	bucketsUseCase := NewBuckets(testObj, settings)
	status := bucketsUseCase.Check("admin", domain.Login)
	if status {
		t.Log("status must be false because markers are over")
	}
	bucket.Marker = 10
	status = bucketsUseCase.Check("admin", domain.Login)
	if !status {
		t.Log("status must be true because markers are enough")
	}
	if bucket.Marker != 9 {
		t.Log("marker must be 9 because subtract 1 from 10")
	}


	status = bucketsUseCase.Check("admin", domain.Login)




}
