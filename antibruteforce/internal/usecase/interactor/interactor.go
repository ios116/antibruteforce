package interactor

import (
	"antibruteforce/internal/domain/entities"
	"antibruteforce/internal/domain/exceptions"
	"context"
	"net"
)

// ConnectorUseCase interface to interaction between use cases
type ConnectorUseCase interface {
	CheckRequest(request *entities.Request) (bool, error)
}

// SubnetChecker check if net is black or white list
type SubnetChecker interface {
	CheckSubnet(ctx context.Context, ip *net.IPNet) (entities.IPKind, error)
}

// BucketsChecker checking bucket or creating if not exist
type BucketsChecker interface {
	CheckOrCreateBucket(request string, kind entities.KindBucket) (bool, error)
}

// Connector interaction between use cases
type Connector struct {
	IP     SubnetChecker
	Bucket BucketsChecker
}

// NewConnector constructor
func NewConnector(IP SubnetChecker, bucket BucketsChecker) *Connector {
	return &Connector{IP: IP, Bucket: bucket}
}

// CheckRequest checking a request by ip login and passport
func (i *Connector) CheckRequest(request *entities.Request) (bool, error) {
	if err := request.Validation(); err != nil {
		return false, err
	}
	ip := net.ParseIP(request.IP)
	if ip == nil {
		return false, exceptions.IPRequired
	}
	ctx := context.Background()
	ipNet := &net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(32, 32),
	}
	kind, err := i.IP.CheckSubnet(ctx, ipNet)
	if err != nil {
		return false, err
	}
	switch kind {
	case entities.Black:
		return false, exceptions.IPInBlackList
	case entities.White:
		return true, nil
	}

	var mainErr error
	IPStatus, err := i.Bucket.CheckOrCreateBucket(request.IP, entities.IP)
	if err != nil {
		mainErr = err
	}
	loginStatus, err := i.Bucket.CheckOrCreateBucket(request.Login, entities.Login)
	if err != nil {
		mainErr = err
	}
	passwordStatus, err := i.Bucket.CheckOrCreateBucket(request.Password, entities.Password)
	if err != nil {
		mainErr = err
	}
	status := IPStatus && loginStatus && passwordStatus
	return status, mainErr
}
