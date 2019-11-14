package entities

import (
	"context"
	"net"
	"time"
)

type KindIp string

const (
	Black KindIp = "black"
	White        = "white"
)

type IPList struct {
	ID          int64
	Kind        KindIp
	IP          *net.IPNet
	DateCreated time.Time
}

type IPManager interface {
	Add(ctx context.Context, list *IPList) error
	DeleteByIp(ctx context.Context, ip *net.IPNet) error
	GetByIP(ctx context.Context, ip *net.IPNet) (*IPList, error)
}
