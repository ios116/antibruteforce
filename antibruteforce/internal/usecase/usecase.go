package usecase

import (
	"antibruteforce/internal/config"
	"antibruteforce/internal/domain/entities"
	"antibruteforce/internal/domain/exceptions"
	"context"
	"net"
	"time"
)

// BucketsManager интерфейс позводляющий проверить наличие свободных маркеров и удалить устаревший bucket
type BucketsManager interface {
	// for buckets
	GetBucketByHash(hash *entities.Hash) (*entities.Bucket, error)
	CreateBucket(hash *entities.Hash) (*entities.Bucket, error)
	CheckBucket(bucket *entities.Bucket) (bool, error)
	BucketCollector(ctx context.Context)
    // for ip
	AddIpToList(ctx context.Context, ip *entities.IPItem) error
	DeleteByIP(ctx context.Context, ip *net.IPNet) error
	GetByIP(ctx context.Context, ip *net.IPNet) (*entities.IPItem, error)
    // for request
	CheckRequest(request entities.Request) (bool, error)
}

// BucketUseCase содержет хранилище buckets, настройки для разного bucket type и канал для удаления неиспользуемых buckets по таймауту
type BucketUseCase struct {
	BucketStore entities.BucketStoreManager
	IPStore     entities.IPStoreManager
	Settings    *config.Settings
	Callback    chan *entities.Hash
}

// NewBuckets создание экземпляра buckets
func NewBuckets(store entities.BucketStoreManager, settings *config.Settings) *BucketUseCase {
	callback := make(chan *entities.Hash)
	return &BucketUseCase{BucketStore: store, Settings: settings, Callback: callback}
}

// GetBucketByHash get bucket by hash
func (b *BucketUseCase) GetBucketByHash(hash *entities.Hash) (*entities.Bucket, error) {
	if hash == nil {
		return nil, exceptions.KeyRequired
	}
	var bucket *entities.Bucket
	bucket, err := b.BucketStore.Get(hash)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

// CreateBucket new bucket with hash
func (b *BucketUseCase) CreateBucket(hash *entities.Hash) (*entities.Bucket, error) {
	if hash == nil {
		return nil, exceptions.KeyRequired
	}
	var bucket *entities.Bucket
	duration := time.Second * time.Duration(b.Settings.Duration)

	switch hash.Kind {
	case entities.Login:
		bucket = entities.NewBucket(b.Settings.LoginRequests, duration, hash, b.Callback)
	case entities.Password:
		bucket = entities.NewBucket(b.Settings.PasswordRequests, duration, hash, b.Callback)
	case entities.IP:
		bucket = entities.NewBucket(b.Settings.IPRequests, duration, hash, b.Callback)
	default:
		return nil, exceptions.TypeNotFound
	}
	err := b.BucketStore.Add(hash, bucket)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

// CheckBucket checks if a limit exist.
func (b *BucketUseCase) CheckBucket(bucket *entities.Bucket) (bool, error) {
	if bucket == nil {
		return false, exceptions.NilValue
	}
	if !bucket.Counter() {
		return false, exceptions.LimitReached
	}
	return true, nil
}

// BucketCollector удаление устаревшего bucket по таймауту, в канал отправляется  bucket's hash
func (b *BucketUseCase) BucketCollector(ctx context.Context) {
	for {
		select {
		case hash := <-b.Callback:
			b.BucketStore.Delete(hash)
		case <-ctx.Done():
			return
		}
	}
}

// AddIpToList adding to ip to list
func (b *BucketUseCase) AddIpToList(ctx context.Context, ip *entities.IPItem) error {
    return b.IPStore.Add(ctx, ip)
}
// DeleteByIP delete ip from list
func (b *BucketUseCase) DeleteByIP(ctx context.Context, ip *net.IPNet) error {
	return b.IPStore.DeleteByIP(ctx, ip)
}

// GetByIP
func (b *BucketUseCase) GetByIP(ctx context.Context, ip *net.IPNet) (*entities.IPItem, error) {
    return b.IPStore.GetByIP(ctx,ip)
}

func (b *BucketUseCase) CheckRequest(request *entities.Request) (bool, error) {
	if err := request.Validation(); err !=nil {
		return false, err
	}
	//ctx := context.Background()
	loginHash := entities.NewHash(entities.Login,request.Login)
	//passwordHash := entities.NewHash(entities.Password,request.Password)
	//loginHash := entities.NewHash(entities.Login,request.IP.String())
	//TODO: check ip from blacklist
    bucket, err := b.GetBucketByHash(loginHash)
    if err == nil {
        bucket, err =b.CreateBucket(loginHash)
        if err != nil {
        	return false, err
		}
	}
    status, err :=b.CheckBucket(bucket)
    return status, err
}
