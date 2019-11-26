package bucketusecase

import (
	"antibruteforce/internal/config"
	"antibruteforce/internal/domain/entities"
	"antibruteforce/internal/domain/exceptions"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// BucketsUseCase интерфейс позводляющий проверить наличие свободных маркеров и удалить устаревший bucket
type BucketsUseCase interface {
	// for buckets
	GetBucketByHash(hash entities.Hash) (*entities.Bucket, error)
	CreateBucket(hash entities.Hash) (*entities.Bucket, error)
	CheckBucket(bucket *entities.Bucket) (bool, error)
	TotalBuckets() int
	BucketCollector()
	ResetBucket(hash entities.Hash) error
	CheckOrCreateBucket(request string, kind entities.KindBucket) (bool, error)
}

// BucketService содержет хранилище buckets, настройки для разного bucket type и канал для удаления неиспользуемых buckets по таймауту
type BucketService struct {
	BucketStore entities.BucketStoreManager
	Settings    *config.Settings
	Callback    chan entities.Hash
	Logger      *zap.Logger
}

// NewBucketService создание экземпляра buckets
func NewBucketService(store entities.BucketStoreManager, settings *config.Settings, logger *zap.Logger) *BucketService {
	callback := make(chan entities.Hash)
	return &BucketService{BucketStore: store, Settings: settings, Callback: callback, Logger: logger}
}

// GetBucketByHash get bucket by hash
func (b *BucketService) GetBucketByHash(hash entities.Hash) (*entities.Bucket, error) {
	if err := hash.Validation(); err != nil {
		return nil, err
	}
	var bucket *entities.Bucket
	bucket, err := b.BucketStore.Get(hash)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

// CreateBucket new bucket with hash
func (b *BucketService) CreateBucket(hash entities.Hash) (*entities.Bucket, error) {
	if err := hash.Validation(); err != nil {
		return nil, err
	}
	var bucket *entities.Bucket
	duration := time.Second * time.Duration(b.Settings.Duration)

	switch hash.Kind {
	case entities.Login:
		bucket = entities.NewBucket(b.Settings.LoginLimit, duration, hash, b.Callback)
	case entities.Password:
		bucket = entities.NewBucket(b.Settings.PasswordLimit, duration, hash, b.Callback)
	case entities.IP:
		bucket = entities.NewBucket(b.Settings.IPLimit, duration, hash, b.Callback)
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
func (b *BucketService) CheckBucket(bucket *entities.Bucket) (bool, error) {
	if bucket == nil {
		return false, fmt.Errorf("CheckBucket: %w", exceptions.NilValue)
	}
	if !bucket.Counter() {
		return false, fmt.Errorf("CheckBucket: %w", exceptions.LimitReached)
	}
	return true, nil
}

// TotalBuckets total buckets in memory
func (b *BucketService) TotalBuckets() int {
	return b.BucketStore.TotalBuckets()
}

//ResetBucket reset bucket by hash
func (b *BucketService) ResetBucket(hash entities.Hash) error {
	if err := b.BucketStore.Delete(hash); err != nil {
		return err
	}
	return nil
}

// BucketCollector удаление устаревшего bucket по таймауту, в канал отправляется  bucket's hash
func (b *BucketService) BucketCollector() {
	for {
		hash := <-b.Callback
		err := b.ResetBucket(hash)
		if err != nil {
			err = fmt.Errorf("bucket collector: %w", err)
			b.Logger.Error(err.Error())
		}
		b.Logger.Info("bucket was deleted", zap.String("hash", hash.Key))
	}
}

// CheckOrCreateBucket check once bucket may be password, login, ip
func (b *BucketService) CheckOrCreateBucket(request string, kind entities.KindBucket) (bool, error) {
	hash := entities.NewHash(kind, request)
	bucket, err := b.GetBucketByHash(hash)
	if bucket == nil {
		bucket, err = b.CreateBucket(hash)
		if err != nil {
			return false, fmt.Errorf("CheckOrCreateBucket: %w", err)
		}
	}
	status, err := b.CheckBucket(bucket)
	if err != nil {
		return false, fmt.Errorf("CheckOrCreateBucket: %w", err)
	}
	return status, nil
}
