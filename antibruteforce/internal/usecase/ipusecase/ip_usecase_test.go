package ipusecase

import (
	"antibruteforce/internal/config"
	"antibruteforce/internal/domain/entities"
	"context"
	"net"
	"testing"
	"time"
)

func TestIpUseCase(t *testing.T) {
	ipv4Addr1, ipv4Net1, _ := net.ParseCIDR("192.168.0.1/24")
	ipv4Addr2, ipv4Net2, _ := net.ParseCIDR("91.225.77.47/32")
	testData := []struct {
		ip            net.IP
		net           *net.IPNet
		expectedIP    []*entities.IPListRow
		expectedError error
		expectedKind  entities.IPKind
	}{
		{
			expectedIP:    []*entities.IPListRow{{ID: 1, Kind: entities.Black, IP: ipv4Net1, DateCreated: time.Now()}},
			ip:            ipv4Addr1,
			net:           ipv4Net1,
			expectedKind:  entities.Black,
			expectedError: nil,
		},
		{
			expectedIP:    []*entities.IPListRow{{}},
			ip:            ipv4Addr2,
			net:           ipv4Net2,
			expectedError: nil,
			expectedKind:  "",
		},
	}

	testObj := new(MockedIPStore)
	settings := config.NewSettings()
	newIPService := NewIPService(settings, testObj)
	ctx := context.Background()

	t.Log(ipv4Addr1, ipv4Net1)
	ip := &entities.IPListRow{
		Kind:        entities.Black,
		IP:          ipv4Net1,
		DateCreated: time.Time{},
	}
	t.Run("Add Ip", func(t *testing.T) {
		testObj.On("Add", ctx, ip).Return(nil)
		err := newIPService.AddNet(ctx, ip)
		if err != nil {
			t.Fatal(err)
		}
		err = newIPService.AddNet(ctx, nil)
		if err == nil {
			t.Fatal(err)
		}
	})

	t.Run("Check ip", func(t *testing.T) {
		for _, item := range testData {
			testObj.On("GetSubnetBySubnet", ctx, item.net).Return(item.expectedIP, nil)
			kind, err := newIPService.CheckSubnet(ctx, item.net)
			if err != item.expectedError {
				t.Fatal(err)
			}
			if kind != item.expectedKind {
				t.Fatalf("expectedKind of list from test data %s but we have %s", item.expectedKind, kind)
			}
		}
	})
}
