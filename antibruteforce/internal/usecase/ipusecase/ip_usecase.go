package ipusecase

import (
	"antibruteforce/internal/config"
	"antibruteforce/internal/domain/entities"
	"antibruteforce/internal/domain/exceptions"
	"context"
	"net"
)

// IPUseCase it's interface for use case ip address
type IPUseCase interface {
	AddNet(ctx context.Context, ip *entities.IPItem) error
	DeleteNet(ctx context.Context, ip *net.IPNet) error
	CheckSubnetAsString(ctx context.Context, address string) (*entities.IPKind, error)
	CheckIpAsString(ctx context.Context, address string) (*entities.IPKind, error)
	checkSubnet(ctx context.Context, ip *net.IPNet) (*entities.IPKind, error)
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
func (b *IPService) AddNet(ctx context.Context, ip *entities.IPItem) error {
	if ip == nil {
		return exceptions.IPRequired
	}
	return b.IPStore.Add(ctx, ip)
}

// DeleteByIP delete ip from list
func (b *IPService) DeleteNet(ctx context.Context, ip *net.IPNet) error {
	if ip == nil {
		return exceptions.IPRequired
	}
	return b.IPStore.DeleteByIP(ctx, ip)
}

// CheckSubnet checks the subnet, whether it is whitelisted or blacklisted
// If it does not contain a return nil
func (b *IPService) CheckSubnet(ctx context.Context, ip *net.IPNet) (entities.IPKind, error) {
	if ip == nil {
		return "",exceptions.IPRequired
	}
	res, err := b.IPStore.GetSubnetBySubnet(ctx, ip)
	if err != nil {
		return "", err
	}
	switch len(res) {
	case 0:
		return "", nil
	default:
		return res[0].Kind, nil
	}
}


