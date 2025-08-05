package drivers

import "errors"

var (
	ErrInvalidConfigStruct = errors.New("invalid configuration structure")
)

var (
	ErrReaderNotExist     = errors.New("reader does not exists")
	ErrReaderUserNotExist = errors.New("reader user does not exists")
)
