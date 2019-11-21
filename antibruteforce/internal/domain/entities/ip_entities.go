package entities

import (
	"context"
	"net"
	"time"
)

// IPKind may by black(IP is blocked) or white(IP is approved)
type IPKind string

const (
	Black IPKind = "black"
	White        = "white"
)

// IPItem ip item
type IPItem struct {
	ID          int64
	Kind        IPKind
	IP          *net.IPNet
	DateCreated time.Time
}

// IPStoreManager is interface for ips storage
type IPStoreManager interface {
	Add(ctx context.Context, list *IPItem) error
	DeleteByIP(ctx context.Context, ip *net.IPNet) error
	GetSubnetBySubnet(ctx context.Context, ip *net.IPNet) ([]*IPItem, error)
}
