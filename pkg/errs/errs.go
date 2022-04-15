package errs

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidInput    = errors.New("invalid input")
	ErrLinkerInternal  = errors.New("proto internal error")
	ErrStorageInternal = errors.New("storage internal error")
	ErrLinkExists      = errors.New("link already saved")
	ErrLinkRemoved     = errors.New("link was removed")
)
