package limiter_manager

import "sync"

// key is service:service_part
type LimiterManagerStorage struct {
	managers map[string]*LimiterManager
	mu       sync.RWMutex
}

func NewLimiterManagerStorage() *LimiterManagerStorage {
	return &LimiterManagerStorage{
		managers: make(map[string]*LimiterManager),
	}
}

func (s *LimiterManagerStorage) Set(key string, manager *LimiterManager) {
	s.mu.Lock()
	s.managers[key] = manager
	s.mu.Unlock()
}

func (s *LimiterManagerStorage) Get(key string) (*LimiterManager, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if manager, ok := s.managers[key]; ok {
		return manager, nil
	}

	return nil, ErrManagerNotFound
}
