package usecase

import (
	"antibruteforce/config"
	"antibruteforce/domain"
	"time"
)

// Manager интерфейс позводляющий проверить наличие свободных маркеров и удалить устаревший bucket
type BucketsManager interface {
	Check(key string, kind domain.Kind) bool
	DeleteBucket(key string)
}

// Buckets содержет хранилище buckets, настройки для разного bucket type и канал для удаления неиспользуемых buckets по таймауту
type Buckets struct {
	Store    domain.StoreManager
	Settings config.Settings
	Callback chan string
}

// NewBuckets создание экземпляра buckets
func NewBuckets(store domain.StoreManager, settings config.Settings) *Buckets {
	callback := make(chan string)
	return &Buckets{Store: store, Settings: settings, Callback: callback}
}

// Check проверка есть ли доступные запросы. Количество доступных запросов зависимости от bucket type
func (b *Buckets) Check(key string, kind domain.Kind) bool {
	bucket := b.Store.Get(key)
	if bucket != nil {
		return bucket.Counter()
	} else {
		duration := time.Second * time.Duration(b.Settings.Duration)
		var bucket *domain.Bucket
		switch kind {
		case domain.Login:
			domain.NewBucket(b.Settings.LoginN, duration, key, b.Callback)
		case domain.Password:
			domain.NewBucket(b.Settings.PasswordM, duration, key, b.Callback)
		case domain.Ip:
			domain.NewBucket(b.Settings.IpK, duration, key, b.Callback)
		}
		b.Store.Add(key, bucket)
		return true
	}
}

// DeleteBucket дуаление устаревшего bucket по таймауту, в канал отправляется  bucket's key
func (b *Buckets) DeleteBucket(key string) {
	for {
		key := <-b.Callback
		b.Store.Delete(key)
	}
}
