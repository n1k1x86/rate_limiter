package bucket_storage

import (
	"rate-limiter/internal/services/limiter"
	"sync"
)

// key is service:service_part:subject
type BucketStorage struct {
	prefixKey string
	buckets   map[string]*limiter.RateLimiter
	mu        sync.Mutex
}

func NewBucketStorage(service, servicePart string) *BucketStorage {
	return &BucketStorage{
		prefixKey: service + ":" + servicePart,
		buckets:   make(map[string]*limiter.RateLimiter),
	}
}

func (s *BucketStorage) key(subject string) string {
	return s.prefixKey + ":" + subject
}

func (s *BucketStorage) GetOrCreate(subject string, create func() (*limiter.RateLimiter, error)) (*limiter.RateLimiter, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.key(subject)

	if l, ok := s.buckets[key]; ok {
		return l, nil
	}

	candidate, err := create()
	if err != nil {
		return nil, err
	}

	s.buckets[key] = candidate

	return candidate, nil
}
