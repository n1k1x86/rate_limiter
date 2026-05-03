package limiter_manager

import (
	"rate-limiter/internal/services/bucket_storage"
	"rate-limiter/internal/services/limiter"
	"sync"
)

type Mode string

const (
	ModeEnforce  Mode = "enforce"
	ModeAllowAll Mode = "allow_all"
	ModeDenyAll  Mode = "deny_all"
)

type LimiterManager struct {
	service       string
	servicePart   string
	bucketStorage *bucket_storage.BucketStorage
	rule          limiter.RateLimiterRule
	mode          Mode
	mu            sync.RWMutex
}

func NewLimiterManager(service, servicePart string, bucketStorage *bucket_storage.BucketStorage, rule *limiter.RateLimiterRule) (*LimiterManager, error) {
	if rule == nil {
		return nil, ErrLimiterRuleIsNil
	}

	if bucketStorage == nil {
		return nil, ErrBucketStorageIsNil
	}

	return &LimiterManager{
		service:       service,
		servicePart:   servicePart,
		bucketStorage: bucketStorage,
		mode:          ModeEnforce,
		rule:          *rule,
	}, nil
}

type Decision struct {
	Allowed     bool   `json:"allowed"`
	Service     string `json:"service"`
	ServicePart string `json:"service_part"`
	Subject     string `json:"subject"`
	Mode        Mode   `json:"mode"`
}

func (m *LimiterManager) snapshot() Mode {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.mode
}

func (m *LimiterManager) SetMode(mode Mode) error {
	if mode != ModeAllowAll && mode != ModeDenyAll && mode != ModeEnforce {
		return ErrModeUndefined
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.mode = mode
	return nil
}

func (m *LimiterManager) decision(subject string, mode Mode, allowed bool) *Decision {
	return &Decision{
		Allowed:     allowed,
		Service:     m.service,
		ServicePart: m.servicePart,
		Subject:     subject,
		Mode:        mode,
	}
}

func (m *LimiterManager) Allow(subject string) (*Decision, error) {

	mode := m.snapshot()

	switch mode {
	case ModeEnforce:
		rateLimiter, err := m.bucketStorage.GetOrCreate(subject, func() (*limiter.RateLimiter, error) {
			l, err := limiter.NewLimiter(&m.rule)
			if err != nil {
				return nil, err
			}
			return l, nil
		})

		if err != nil {
			return nil, err
		}

		if rateLimiter == nil {
			return nil, ErrRateLimiterIsNil
		}

		allow := rateLimiter.Allow()

		return m.decision(subject, mode, allow), nil
	case ModeAllowAll:
		return m.decision(subject, mode, true), nil
	case ModeDenyAll:
		return m.decision(subject, mode, false), nil
	}
	return m.decision(subject, mode, false), ErrModeUndefined
}
