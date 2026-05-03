package limiter

import "errors"

var (
	ErrInvalidRuleParams = errors.New("invalid rule params")
	ErrRuleIsNil         = errors.New("rule is nil")
)
