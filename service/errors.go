package service

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailExists        = errors.New("email already registered")
	ErrNotFound           = errors.New("resource not found")
	ErrForbidden          = errors.New("forbidden")
	ErrZoneFull           = errors.New("parking zone is full")
)
