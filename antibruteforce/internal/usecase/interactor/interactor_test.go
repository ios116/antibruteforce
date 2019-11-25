package interactor

import (
	"antibruteforce/internal/domain/entities"
	"context"
	"github.com/stretchr/testify/mock"
	"net"
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

type SubnetCheckerMock struct {
	mock.Mock
}

func (m *SubnetCheckerMock) CheckSubnet(ctx context.Context, ip *net.IPNet) (entities.IPKind, error) {
	args := m.Called(ctx, ip)
	return args.Get(0).(entities.IPKind), args.Error(1)
}

type BucketsCheckerMock struct {
	mock.Mock
}

func (m *BucketsCheckerMock) GetBucketByHash(hash *entities.Hash) (*entities.Bucket, error) {
	args := m.Called(hash)
	return args.Get(0).(*entities.Bucket), args.Error(1)
}
func (m *BucketsCheckerMock) CreateBucket(hash *entities.Hash) (*entities.Bucket, error) {
	args := m.Called(hash)
	return args.Get(0).(*entities.Bucket), args.Error(1)
}
func (m *BucketsCheckerMock) CheckBucket(bucket *entities.Bucket) (bool, error) {
	args := m.Called(bucket)
	return args.Bool(0), args.Error(1)
}

func TestConnector(t *testing.T) {
	subnetUseCase := new(SubnetCheckerMock)
	bucketUseCase := new(BucketsCheckerMock)
	connector := NewConnector(subnetUseCase, bucketUseCase)

	testData := []struct {
		request string
		kind entities.KindBucket
		expected bool
		err error
	}{
      {"admin", entities.Login, false, nil},
      {"123456", entities.Password, true, nil},
      {"127.0.0.1", entities.IP, true, nil},
	}

    for _, item :=range testData {
		//GetBucketByHash
		//CreateBucket
		// CheckBucket
		hash:=entities.NewHash(item.kind, item.request)
		bucket:=&entities.Bucket{Marker:0,Duration: 0,}
		bucketUseCase.On("GetBucketByHash",hash)
		status, err := connector.CheckOnceBucket(item.request, item.kind)
		assert.Equal(t,status, item.expected)
		assert.Equal(t,err,item.err)

	}
}
