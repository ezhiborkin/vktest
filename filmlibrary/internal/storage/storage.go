package storage

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUserExists    = errors.New("user exists")
	ErrActorNotFound = errors.New("actor not found")
	ErrActorExists   = errors.New("actor exists")
	ErrMovieNotFound = errors.New("movie not found")
	ErrMovieExists   = errors.New("movie exists")
)
