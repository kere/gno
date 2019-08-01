package db

import (
	"errors"
	"time"
)

var (
	// ErrNoField err
	ErrNoField = errors.New("field not found")

	// ErrType err
	ErrType = errors.New("db convert to value type error")

	// EmptyTime empty
	EmptyTime = time.Time{}
)
