package entities

import (
	"context"
	"net"
	"time"
)

// IPKind may by black(IP is blocked) or white(IP is approved)
type IPKind string

//Black, White type of ip list
const (
	Black IPKind = "black"
	White        = "white"
)

// IPListRow ip item
type IPListRow struct {
	ID          int64
	Kind        IPKind
	IP          *net.IPNet
	DateCreated time.Time
}

// IPStoreManager is interface for ips storage
type IPStoreManager interface {
	Add(ctx context.Context, list *IPListRow) error
	DeleteByIP(ctx context.Context, ip *net.IPNet) error
	GetSubnetBySubnet(ctx context.Context, ip *net.IPNet) ([]*IPListRow, error)
}
