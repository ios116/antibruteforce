package bucketusecase

import (
	"antibruteforce/internal/config"
	"antibruteforce/internal/domain/entities"
	"antibruteforce/internal/store/bucketstore"
	"context"
	"os"
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

var bucketService *BucketService

func TestMain(m *testing.M) {
	bucketStore := bucketstore.NewBucketStore()
	// Create settings
	settings := config.NewSettings()
	// set 1 request in 3 seconds to login
	settings.Duration = 3
	settings.LoginLimit = 1
	// Create buckets use case
	bucketService = NewBucketService(bucketStore, settings)
	code := m.Run()
	os.Exit(code)
}

func TestGet(t *testing.T) {

	var bucket *entities.Bucket
	var err error

	hash := entities.NewHash(entities.Login, "admin")

	t.Run("GetBucketByHash bucket if not exist", func(t *testing.T) {
		_, err = bucketService.GetBucketByHash(hash)
		if err == nil {
			t.Fatal("bucket must be nil and error should be")
		}
	})

	t.Run("CreateBucket bucket with login type and set value admin to login", func(t *testing.T) {
		_, err = bucketService.CreateBucket(hash)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Checking the presence of a bucket after adding", func(t *testing.T) {
		bucket, err = bucketService.GetBucketByHash(hash)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("CheckBucket available request if enough markers", func(t *testing.T) {
		status, err := bucketService.CheckBucket(bucket)
		if err != nil {
			t.Fatal(err)
		}
		if !status {
			t.Fatal("Status must be true because there are enough markers")
		}
	})

	t.Run("CheckBucket available request if not enough markers", func(t *testing.T) {
		status, err := bucketService.CheckBucket(bucket)
		if err == nil {
			t.Fatal(err)
		}
		if status {
			t.Fatal("Status must be false because there are not enough markers")
		}
	})

	// Running collector of buckets
	ctx := context.Background()
	go bucketService.BucketCollector(ctx)
	time.Sleep(time.Second * 4)
	t.Run("CheckBucket for bucket removal after the expiration of a lifetime", func(t *testing.T) {
		_, err = bucketService.GetBucketByHash(hash)
		if err == nil {
			t.Fatal("bucket must be nil")
		}
	})
}

//func BenchmarkBucket(b *testing.B) {
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			hash := &entities.Hash{
//				Kind: entities.Login,
//				Key:  "admin",
//			}
//			bucketService.CreateBucket(hash)
//		}
//	})
//}

func TestMem(t *testing.T) {
	requests := []struct {
		kind entities.KindBucket
		key  string
	}{
		{kind: entities.Login, key: "admin"},
		{kind: entities.Login, key: "manager"},
		{kind: entities.Login, key: "user"},
		{kind: entities.Login, key: "another_user"},
	}
	for _, request := range requests {
		hash := &entities.Hash{
			Kind: request.kind,
			Key:  request.key,
		}
		bucketService.CreateBucket(hash)
	}
	t.Log("Total=", bucketService.TotalBuckets())
	time.Sleep(2 * time.Second)
	t.Log("Total=", bucketService.TotalBuckets())

}
