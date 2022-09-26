package merge

import (
	"errors"
)

var (
	ErrKindNotSupported = errors.New("must be a struct, map or slice")
	ErrNotAdrressable   = errors.New("must not be unadrressable value")
	ErrNilValue         = errors.New("must not be nil value")
	ErrInvalidValue     = errors.New("must not be invalid value")
	ErrUnknownRange     = errors.New("unknown range")
	ErrUnknownResolver  = errors.New("unknown resolver")
	ErrNotSettable      = errors.New("must be settable")
)
