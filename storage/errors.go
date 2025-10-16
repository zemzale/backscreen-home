package storage

import "errors"

var (
	ErrDuplicate = errors.New("duplicate entry")
	ErrNotFound  = errors.New("not found")
)
