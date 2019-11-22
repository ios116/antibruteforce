package ipusecase

import (
	"antibruteforce/internal/config"
	"antibruteforce/internal/domain/entities"
	"context"
	"github.com/stretchr/testify/mock"
	"net"
	"testing"
	"time"
)

type MockedIPStore struct {
	mock.Mock
}

func (m MockedIPStore) Add(ctx context.Context, list *entities.IPItem) error {
	args := m.Called(ctx, list)
	return args.Error(0)
}
func (m MockedIPStore) DeleteByIP(ctx context.Context, ip *net.IPNet) error {
	args := m.Called(ctx, ip)
	return args.Error(0)
}
func (m MockedIPStore) GetSubnetBySubnet(ctx context.Context, ip *net.IPNet) ([]*entities.IPItem, error) {
	args := m.Called(ctx, ip)
	return args.Get(0).([]*entities.IPItem), args.Error(1)
}

func TestIpUseCase(t *testing.T) {
	ipv4Addr1, ipv4Net1, _ := net.ParseCIDR("192.168.0.1/24")
	ipv4Addr2, ipv4Net2, _ := net.ParseCIDR("91.225.77.47/32")

	testData := []struct {
		ip            net.IP
		net           *net.IPNet
		expectedIP    []*entities.IPItem
		expectedError error
		expectedKind  entities.IPKind
	}{
		{
			expectedIP: []*entities.IPItem{
				&entities.IPItem{
					ID:          1,
					Kind:        entities.Black,
					IP:          ipv4Net1,
					DateCreated: time.Now(),
				},
			},
			ip:            ipv4Addr1,
			net:           ipv4Net1,
			expectedKind:  entities.Black,
			expectedError: nil,
		},
		{
			expectedIP: []*entities.IPItem{
				&entities.IPItem{},
			},
			ip:            ipv4Addr2,
			net:           ipv4Net2,
			expectedError: nil,
			expectedKind:  "",
		},
	}

	testObj := new(MockedIPStore)
	settings := config.NewSettings()
	newIpService := NewIPService(settings, testObj)
	ctx := context.Background()

	t.Log(ipv4Addr1, ipv4Net1)
	ip := &entities.IPItem{
		Kind:        entities.Black,
		IP:          ipv4Net1,
		DateCreated: time.Time{},
	}
	t.Run("Add Ip", func(t *testing.T) {
		testObj.On("Add", ctx, ip).Return(nil)
		err := newIpService.AddNet(ctx, ip)
		if err != nil {
			t.Fatal(err)
		}
		err = newIpService.AddNet(ctx, nil)
		if err == nil {
			t.Fatal(err)
		}
	})

	t.Run("Check ip", func(t *testing.T) {
		for _, item := range testData {
			testObj.On("GetSubnetBySubnet", ctx, item.net).Return(item.expectedIP, nil)
			kind, err := newIpService.CheckSubnet(ctx, item.net)
			if err != item.expectedError {
				t.Fatal(err)
			}
			if kind != item.expectedKind {
				t.Fatalf("expectedKind of list from test data %s but we have %s", item.expectedKind, kind, )
			}
		}
	})
}



