package usecase

import "antibruteforce/internal/domain/entities"

// InteractorUseCase 
type InteractorUseCase interface {
	CheckRequest(request *entities.Request) (bool, error)
	CheckOnceBucket(request string, kind entities.KindBucket) (bool, error)
}

// Interactor 
type Interactor struct {
	IP     IPService
	Bucket BucketService
}

func NewInteractor(IP IPService, bucket BucketService) *Interactor {
	return &Interactor{IP: IP, Bucket: bucket}
}

// CheckRequest
func (b *Interactor) CheckRequest(request *entities.Request) (bool, error) {
	if err := request.Validation(); err != nil {
		return false, err
	}
	//ctx := context.Background()
	//passwordHash := entities.NewHash(entities.Password,request.Password)
	IPStatus, err := b.CheckOnceBucket(request.IP, entities.IP)
	if err != nil {
		return false, err
	}
	loginStatus, err := b.CheckOnceBucket(request.IP, entities.Login)
	if err != nil {
		return false, err
	}
	passwordStatus, err := b.CheckOnceBucket(request.IP, entities.Password)
	if err != nil {
		return false, err
	}
	status := IPStatus && loginStatus && passwordStatus
	return status, nil
}

// CheckOnceBucket
func (b *Interactor) CheckOnceBucket(request string, kind entities.KindBucket) (bool, error) {
	hash := entities.NewHash(kind, request)
	bucket, err := b.Bucket.GetBucketByHash(hash)
	if err == nil {
		bucket, err = b.Bucket.CreateBucket(hash)
		if err != nil {
			return false, err
		}
	}
	status, err := b.Bucket.CheckBucket(bucket)
	return status, err
}
