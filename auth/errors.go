package auth

import "errors"

var (
	ErrNoAccounts      = errors.New("no accounts configured")
	ErrAccountNotFound = errors.New("account not found")
)
