package usecase

import (
	"antibruteforce/config"
	"antibruteforce/domain/entities"
	"antibruteforce/domain/exceptions"
	"context"
	"time"
)

// Manager интерфейс позводляющий проверить наличие свободных маркеров и удалить устаревший bucket
type BucketsManager interface {
	Check(key string, kind entities.Kind) bool
	DeleteBucket(key string)
}

// Buckets содержет хранилище buckets, настройки для разного bucket type и канал для удаления неиспользуемых buckets по таймауту
type Buckets struct {
	Store    entities.StoreManager
	Settings *config.Settings
	Callback chan *entities.Hash
}

// NewBuckets создание экземпляра buckets
func NewBuckets(store entities.StoreManager, settings *config.Settings) *Buckets {
	callback := make(chan *entities.Hash)
	return &Buckets{Store: store, Settings: settings, Callback: callback}
}

func (b *Buckets) Get(hash *entities.Hash) (*entities.Bucket, error) {
	if hash == nil {
		return nil, exceptions.KeyRequired
	}
	var bucket *entities.Bucket
	bucket, err := b.Store.Get(hash)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func (b *Buckets) Add(hash *entities.Hash) (*entities.Bucket, error) {
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
	case entities.Ip:
		bucket = entities.NewBucket(b.Settings.IpRequests, duration, hash, b.Callback)
	default:
		return nil, exceptions.TypeNotFound
	}
	err := b.Store.Add(hash, bucket)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

// Check проверка есть ли доступные markers.
func (b *Buckets) Check(bucket *entities.Bucket) (bool, error) {
	if bucket == nil {
		return false, exceptions.BucketsNil
	}
	return bucket.Counter(), nil
}

// DeleteBucket дуаление устаревшего bucket по таймауту, в канал отправляется  bucket's key
func (b *Buckets) BucketCollector(ctx context.Context) {
	for {
		select {
		case hash := <-b.Callback:
			b.Store.Delete(hash)
		case <-ctx.Done():
			return
		}
	}
}
