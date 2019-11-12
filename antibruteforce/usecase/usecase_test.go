package usecase

import (
	"antibruteforce/config"
	"antibruteforce/domain/entities"
	"antibruteforce/domain/exceptions"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockedBucketStore struct {
	mock.Mock
}

func (s *MockedBucketStore) Add(key string, bucket *entities.Bucket) error {
	args := s.Called(key, bucket)
	return args.Error(0)
}
func (s *MockedBucketStore) Delete(key string) error {
	args := s.Called(key)
	return args.Error(0)
}
func (s *MockedBucketStore) Get(key string) (*entities.Bucket, error) {
	args := s.Called(key)
	return args.Get(0).(*entities.Bucket), args.Error(1)
}

func TestGet(t *testing.T) {

	testObj := new(MockedBucketStore)
	settings := config.NewSettings()
	bucketsUseCase := NewBuckets(testObj, settings)
	t.Run("Get bucket", func(t *testing.T) {
		testObj.On("Get", "admin").Return(&entities.Bucket{}, exceptions.BucketsNil)
		_, err := bucketsUseCase.Get("admin")
		if err == nil {
			t.Fatal("bucket not created yet")
		}
	})

	t.Run("Created bucket", func(t *testing.T) {
		duration := time.Second * time.Duration(bucketsUseCase.Settings.Duration)
		bucket := entities.NewBucket(bucketsUseCase.Settings.LoginN, duration, "admin", bucketsUseCase.Callback)
        testObj.On("Add","admin",bucket).Return(nil)
        bucket,err:= bucketsUseCase.Add("admin",entities.Login)
        t.Log(bucket, err)
	})

	// status, _:=bucketsUseCase.Check(bucket)

	//if !status {
	//	t.Log("status must be true because markers are enough")
	//}

	//if status {
	//	t.Log("status must be false because markers are over")
	//}
	//
	//
	//testObj.On("Get", "admin").Return(bucket)
	//
	//status := bucketsUseCase.Check("admin", entities.Login)
	//
	//bucket.Marker = 10
	//status = bucketsUseCase.Check("admin", entities.Login)
	//
	//if bucket.Marker != 9 {
	//	t.Log("marker must be 9 because subtract 1 from 10")
	//}
	//
	//
	//status = bucketsUseCase.Check("admin", entities.Login)
	//
	//

}
