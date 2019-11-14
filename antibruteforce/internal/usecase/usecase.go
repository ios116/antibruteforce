package usecase

import (
	"antibruteforce/internal/config"
	"antibruteforce/internal/domain/entities"
	"antibruteforce/internal/domain/exceptions"
	"context"
	"time"
)

// BucketsManager интерфейс позводляющий проверить наличие свободных маркеров и удалить устаревший bucket
type BucketsManager interface {
	GetBucketByHash(hash *entities.Hash) (*entities.Bucket, error)
	CreateBucket(hash *entities.Hash) (*entities.Bucket, error)
	CheckBucket(bucket *entities.Bucket) (bool, error)
	BucketCollector(ctx context.Context)
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
	return bucket.Counter(), nil
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
