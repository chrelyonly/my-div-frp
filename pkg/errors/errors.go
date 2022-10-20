package errors

import (
	"errors"
)

var (
	ErrMsgType   = errors.New("message type error")
	ErrCtlClosed = errors.New("control is closed")
)
