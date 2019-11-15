package usecase

import (
	"antibruteforce/internal/config"
	"antibruteforce/internal/domain/entities"
	"antibruteforce/internal/store/bucketstore"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
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
	//t.Run("GetBucketByHash bucket", func(t *testing.T) {
	//	testObj.On("GetBucketByHash", "admin").Return(&entities.Bucket{}, exceptions.BucketsNil)
	//	_, err := bucketsUseCase.GetBucketByHash("admin")
	//	if err == nil {
	//		t.Fatal("bucket not created yet")
	//	}
	//})
	//
	//t.Run("Created bucket", func(t *testing.T) {
	//	duration := time.Second * time.Duration(bucketsUseCase.Settings.Duration)
	//	bucket := entities.NewBucket(bucketsUseCase.Settings.LoginRequests, duration, "admin", bucketsUseCase.Callback)
	//    testObj.On("CreateBucket","admin",bucket).Return(nil)
	//    bucket,err:= bucketsUseCase.CreateBucket("admin",entities.Login)
	//    t.Log(bucket, err)
	//})

	// Create bucket store
	bucketStore := bucketstore.NewBucketStore()
	// Create settings
	settings := config.NewSettings()
	// set 1 request in 3 seconds to login
	settings.Duration = 3
	settings.LoginRequests = 1
	// Create buckets use case
	bucketsUseCase := NewBuckets(bucketStore, settings)
	var bucket *entities.Bucket
	var err error

	hash := entities.NewHash(entities.Login, "admin")

	t.Run("GetBucketByHash bucket if not exist", func(t *testing.T) {
		_, err = bucketsUseCase.GetBucketByHash(hash)
		if err == nil {
			t.Fatal("bucket must be nil")
		}
	})

	t.Run("CreateBucket bucket with login type and set value admin to login", func(t *testing.T) {
		_, err = bucketsUseCase.CreateBucket(hash)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Checking the presence of a bucket after adding", func(t *testing.T) {
		bucket, err = bucketsUseCase.GetBucketByHash(hash)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("CheckBucket available request if enough markers", func(t *testing.T) {
		status, err := bucketsUseCase.CheckBucket(bucket)
		if err != nil {
			t.Fatal(err)
		}
		if !status {
			t.Fatal("Status must be true because there are enough markers")
		}
	})

	t.Run("CheckBucket available request if not enough markers", func(t *testing.T) {
		status, err := bucketsUseCase.CheckBucket(bucket)
		if err != nil {
			t.Fatal(err)
		}
		if status {
			t.Fatal("Status must be false because there are not enough markers")
		}
	})

	// Running collector of buckets
	ctx := context.Background()
	go bucketsUseCase.BucketCollector(ctx)
	time.Sleep(time.Second * 4)
	t.Run("CheckBucket for bucket removal after the expiration of a lifetime", func(t *testing.T) {
		_, err = bucketsUseCase.GetBucketByHash(hash)
		if err == nil {
			t.Fatal("bucket must be nil")
		}
	})
}
