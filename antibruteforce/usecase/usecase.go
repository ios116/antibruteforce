package usecase

import (
	"antibruteforce/config"
	"antibruteforce/domain/entities"
	"antibruteforce/domain/exceptions"
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
	Callback chan string
}

// NewBuckets создание экземпляра buckets
func NewBuckets(store entities.StoreManager, settings *config.Settings) *Buckets {
	callback := make(chan string)
	return &Buckets{Store: store, Settings: settings, Callback: callback}
}

func (b *Buckets) Get(key string) (*entities.Bucket, error) {
	if key == "" {
		return nil, exceptions.KeyRequired
	}
	var bucket *entities.Bucket
	bucket, err := b.Store.Get(key)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func (b *Buckets) Add(key string, kind entities.Kind) (*entities.Bucket, error) {
	if key == "" {
		return nil, exceptions.KeyRequired
	}
	var bucket *entities.Bucket
	duration := time.Second * time.Duration(b.Settings.Duration)
	switch kind {
	case entities.Login:
		bucket = entities.NewBucket(b.Settings.LoginN, duration, key, b.Callback)
	case entities.Password:
		bucket = entities.NewBucket(b.Settings.PasswordM, duration, key, b.Callback)
	case entities.Ip:
		bucket = entities.NewBucket(b.Settings.IpK, duration, key, b.Callback)
	default:
		return nil, exceptions.TypeNotFound
	}
	err := b.Store.Add(key, bucket)
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
func (b *Buckets) BucketCollector(key string) {
	for {
		key := <-b.Callback
		b.Store.Delete(key)
	}
}
