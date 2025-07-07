package svcErr

import "fmt"

var (
	ErrInvalidID = fmt.Errorf("invalid id format")
)
