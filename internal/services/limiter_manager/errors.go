package limiter_manager

import "errors"

var (
	ErrManagerNotFound    = errors.New("manager not found error")
	ErrModeUndefined      = errors.New("mode is undefined")
	ErrLimiterRuleIsNil   = errors.New("limiter rules is nil")
	ErrRateLimiterIsNil   = errors.New("rate limiter is nil")
	ErrBucketStorageIsNil = errors.New("bucket storage is nil")
)
