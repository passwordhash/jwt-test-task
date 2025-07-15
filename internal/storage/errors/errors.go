package repoerr

import (
	"errors"
)

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)
