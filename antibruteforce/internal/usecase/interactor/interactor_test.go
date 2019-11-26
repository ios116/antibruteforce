package interactor

import (
	"antibruteforce/internal/domain/entities"
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockedSubnetChecker struct {
	mock.Mock
}

func (s *MockedSubnetChecker) CheckSubnet(ctx context.Context, ip *net.IPNet) (entities.IPKind, error) {
	args := s.Called(ctx, ip)
	return args.Get(0).(entities.IPKind), args.Error(1)
}

type MockedBucketsChecker struct {
	mock.Mock
}

func (s *MockedBucketsChecker) CheckOrCreateBucket(request string, kind entities.KindBucket) (bool, error) {
	args := s.Called(request, kind)
	return args.Bool(0), args.Error(1)
}
func TestConnector(t *testing.T) {
	subnetChecker := new(MockedSubnetChecker)
	bucketChecker := new(MockedBucketsChecker)
	connector := NewConnector(subnetChecker, bucketChecker)

	dataSet := []struct {
		request *entities.Request
		ok      bool
		err     error
		IPKind  entities.IPKind
	}{
		{request: &entities.Request{IP: "127.0.0.1", Login: "admin1", Password: "1"}, ok: true, err: nil, IPKind: entities.White},  // if ip in white list
		{request: &entities.Request{IP: "127.0.0.2", Login: "admin2", Password: "2"}, ok: false, err: nil, IPKind: entities.Black}, // if ip in black list
		{request: &entities.Request{IP: "127.0.0.3", Login: "admin3", Password: "3"}, ok: true, err: nil, IPKind: ""},              // if ip in not list
	}
	ctx := context.Background()
	for _, item := range dataSet {
		ipNet := &net.IPNet{
			IP:   net.ParseIP(item.request.IP),
			Mask: net.CIDRMask(32, 32),
		}
		subnetChecker.On("CheckSubnet", ctx, ipNet).Return(item.IPKind, item.err)
		bucketChecker.On("CheckOrCreateBucket", item.request.IP, entities.IP).Return(true, nil)
		bucketChecker.On("CheckOrCreateBucket", item.request.Login, entities.Login).Return(true, nil)
		bucketChecker.On("CheckOrCreateBucket", item.request.Password, entities.Password).Return(true, nil)
		ok, err := connector.CheckRequest(item.request)
		if ok != item.ok {
			t.Fatalf("status not equal")
		}
		if err != item.err {
			t.Fatalf("err not equal")
		}
	}
}
