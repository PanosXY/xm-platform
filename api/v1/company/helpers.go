package company

import "errors"

var (
	ErrDuplicateKey = errors.New("Duplicate key value violates unique constraint")
)
