package merge

import "errors"

var (
	ErrKindNotSupported = errors.New("kind not supported")
	ErrNotAdrressable   = errors.New("must not be unadrressable value")
	ErrNilValue         = errors.New("must not be nil value")
	ErrInvalidValue     = errors.New("must not be invalid value")
	ErrUnknownRange     = errors.New("unknown range")
	ErrUnknownResolver  = errors.New("unknown resolver")
	ErrNotSettable      = errors.New("must be settable")
	ErrInvalidStrategy  = errors.New("invalid strategy")
)
