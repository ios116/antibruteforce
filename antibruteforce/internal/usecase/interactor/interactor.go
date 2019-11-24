package interactor

import (
	"antibruteforce/internal/domain/entities"
	"antibruteforce/internal/domain/exceptions"
	"antibruteforce/internal/usecase/bucketusecase"
	"antibruteforce/internal/usecase/ipusecase"
	"context"
	"net"
)

// InteractorUseCase interface to interaction between use cases
type InteractorUseCase interface {
	CheckRequest(request *entities.Request) (bool, error)
	CheckOnceBucket(request string, kind entities.KindBucket) (bool, error)
}

// Interactor interaction between use cases
type Interactor struct {
	IP     ipusecase.IPUseCase
	Bucket bucketusecase.BucketsUseCase
}

// NewInteractor constructor
func NewInteractor(IP ipusecase.IPUseCase, bucket bucketusecase.BucketsUseCase) *Interactor {
	return &Interactor{IP: IP, Bucket: bucket}
}

// CheckRequest checking a request by ip login and passport
func (i *Interactor) CheckRequest(request *entities.Request) (bool, error) {
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
		return false, nil
	case entities.White:
		return true, nil
	}

	IPStatus, err := i.CheckOnceBucket(request.IP, entities.IP)
	if err != nil {
		return false, err
	}
	loginStatus, err := i.CheckOnceBucket(request.Login, entities.Login)
	if err != nil {
		return false, err
	}
	passwordStatus, err := i.CheckOnceBucket(request.Password, entities.Password)
	if err != nil {
		return false, err
	}
	status := IPStatus && loginStatus && passwordStatus
	return status, nil
}

// CheckOnceBucket check once bucket may be password, login, ip
func (i *Interactor) CheckOnceBucket(request string, kind entities.KindBucket) (bool, error) {
	hash := entities.NewHash(kind, request)
	bucket, err := i.Bucket.GetBucketByHash(hash)
	if err == nil {
		bucket, err = i.Bucket.CreateBucket(hash)
		if err != nil {
			return false, err
		}
	}
	status, err := i.Bucket.CheckBucket(bucket)
	return status, err
}
