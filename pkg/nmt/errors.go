package nmt

import "errors"

var (
	ErrNilHashFunction  error = errors.New("nil hash function")
	ErrInvalidLeafLen   error = errors.New("invalid leaf length")
	ErrInvalidOrder     error = errors.New("invalid order")
	ErrInvalidNamespace error = errors.New("invalid namespace")
	ErrInvalidRange     error = errors.New("invalid range")
	ErrInvalidLevel     error = errors.New("invalid level")
)
