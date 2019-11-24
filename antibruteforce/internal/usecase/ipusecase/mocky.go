package ipusecase

import (
	"antibruteforce/internal/domain/entities"
	"context"
	"net"

	"github.com/stretchr/testify/mock"
)

// MockedIPStore mocked
type MockedIPStore struct {
	mock.Mock
}

//Add mocked
func (m *MockedIPStore) Add(ctx context.Context, list *entities.IPListRow) error {
	args := m.Called(ctx, list)
	return args.Error(0)
}

//DeleteByIP mocked
func (m *MockedIPStore) DeleteByIP(ctx context.Context, ip *net.IPNet) error {
	args := m.Called(ctx, ip)
	return args.Error(0)
}

//GetSubnetBySubnet mocked
func (m *MockedIPStore) GetSubnetBySubnet(ctx context.Context, ip *net.IPNet) ([]*entities.IPListRow, error) {
	args := m.Called(ctx, ip)
	return args.Get(0).([]*entities.IPListRow), args.Error(1)
}

//GetSubnet mocked
func (m *MockedIPStore) GetSubnet(ctx context.Context, subnet string) ([]*entities.IPListRow, error) {
	args := m.Called(ctx, subnet)
	return args.Get(0).([]*entities.IPListRow), args.Error(1)
}
