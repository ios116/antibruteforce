package ipusecase

import (
	"antibruteforce/internal/domain/entities"
	"context"
	"github.com/stretchr/testify/mock"
	"net"
)

type MockedIPStore struct {
	mock.Mock
}

func (m *MockedIPStore) Add(ctx context.Context, list *entities.IPListRow) error {
	args := m.Called(ctx, list)
	return args.Error(0)
}
func (m *MockedIPStore) DeleteByIP(ctx context.Context, ip *net.IPNet) error {
	args := m.Called(ctx, ip)
	return args.Error(0)
}
func (m *MockedIPStore) GetSubnetBySubnet(ctx context.Context, ip *net.IPNet) ([]*entities.IPListRow, error) {
	args := m.Called(ctx, ip)
	return args.Get(0).([]*entities.IPListRow), args.Error(1)
}
