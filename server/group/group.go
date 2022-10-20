package group

import (
	"errors"
)

var (
	ErrGroupAuthFailed    = errors.New("group auth failed")
	ErrGroupParamsInvalid = errors.New("group params invalid")
	ErrListenerClosed     = errors.New("group listener closed")
	ErrGroupDifferentPort = errors.New("group should have same remote port")
	ErrProxyRepeated      = errors.New("group proxy repeated")
)
