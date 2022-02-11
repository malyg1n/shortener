package errs

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidInput    = errors.New("invalid input")
	ErrLinkerInternal  = errors.New("linker internal error")
	ErrStorageInternal = errors.New("storage internal error")
	ErrLinkExists      = errors.New("link already saved")
)
