package errors

import "errors"

var (
	ErrInvalidOutputType = errors.New("this output type isn't handled")
)
