package usecase

import (
	"antibruteforce/config"
	"antibruteforce/domain/entities"
	"antibruteforce/store"
	"context"
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

	//testObj := new(MockedBucketStore)
	//settings := config.NewSettings()
	//bucketsUseCase := NewBuckets(testObj, settings)
	//t.Run("Get bucket", func(t *testing.T) {
	//	testObj.On("Get", "admin").Return(&entities.Bucket{}, exceptions.BucketsNil)
	//	_, err := bucketsUseCase.Get("admin")
	//	if err == nil {
	//		t.Fatal("bucket not created yet")
	//	}
	//})
	//
	//t.Run("Created bucket", func(t *testing.T) {
	//	duration := time.Second * time.Duration(bucketsUseCase.Settings.Duration)
	//	bucket := entities.NewBucket(bucketsUseCase.Settings.LoginN, duration, "admin", bucketsUseCase.Callback)
	//    testObj.On("Add","admin",bucket).Return(nil)
	//    bucket,err:= bucketsUseCase.Add("admin",entities.Login)
	//    t.Log(bucket, err)
	//})

	// Create bucket store
	bucketStore := store.NewBucketStore()
	// Create settings
	settings := config.NewSettings()
	// set 1 request in 3 seconds to login
	settings.Duration = 3
	settings.LoginN = 1
	// Create buckets use case
	bucketsUseCase := NewBuckets(bucketStore, settings)
	var bucket *entities.Bucket
	var err error

	t.Run("Get bucket if not exist", func(t *testing.T) {
		_, err = bucketsUseCase.Get("login")
		if err == nil {
			t.Fatal("bucket must be nil")
		}
	})

	t.Run("Add bucket with login type and set value admin to login", func(t *testing.T) {
		_, err = bucketsUseCase.Add("admin", entities.Login)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Checking the presence of a bucket after adding", func(t *testing.T) {
		bucket, err = bucketsUseCase.Get("admin")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Check available request if enough markers", func(t *testing.T) {
		status, err := bucketsUseCase.Check(bucket)
		if err != nil {
			t.Fatal(err)
		}
		if !status {
			t.Fatal("Status must be true because there are enough markers")
		}
	})

	t.Run("Check available request if not enough markers", func(t *testing.T) {
		status, err := bucketsUseCase.Check(bucket)
		if err != nil {
			t.Fatal(err)
		}
		if status {
			t.Fatal("Status must be false because there are not enough markers")
		}
	})

	// Running collector of buckets
	ctx:=context.Background()
	go bucketsUseCase.BucketCollector(ctx)
	time.Sleep(time.Second * 4)
	t.Run("Check for bucket removal after the expiration of a lifetime", func(t *testing.T) {
		_, err = bucketsUseCase.Get("admin")
		if err == nil {
			t.Fatal("bucket must be nil")
		}
	})
}
