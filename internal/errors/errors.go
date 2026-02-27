package errors

import "errors"

var (
	ErrInvalidOutputType = errors.New("This output type isn't handled.")
)
