package ipusecase

import (
	"antibruteforce/internal/config"
	"antibruteforce/internal/domain/entities"
	"context"
	"net"
)

type IPUseCase interface {
	AddIpToList(ctx context.Context, ip *entities.IPItem) error
	DeleteByIP(ctx context.Context, ip *net.IPNet) error
	GetByIP(ctx context.Context, ip *net.IPNet) (*entities.IPItem, error)
}

// IPService
type IPService struct {
	Settings *config.Settings
	IPStore  entities.IPStoreManager
}

// NewIPService
func NewIPService(settings *config.Settings, IPStore entities.IPStoreManager) *IPService {
	return &IPService{Settings: settings, IPStore: IPStore}
}

// AddIpToList adding to ip to list
func (b *IPService) AddIpToList(ctx context.Context, ip *entities.IPItem) error {
	return b.IPStore.Add(ctx, ip)
}

// DeleteByIP delete ip from list
func (b *IPService) DeleteByIP(ctx context.Context, ip *net.IPNet) error {
	return b.IPStore.DeleteByIP(ctx, ip)
}

// GetByIP
func (b *IPService) GetByIP(ctx context.Context, ip *net.IPNet) (*entities.IPItem, error) {
	return b.IPStore.GetByIP(ctx, ip)
}