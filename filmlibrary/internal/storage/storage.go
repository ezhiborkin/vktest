package storage

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrMovieNotFound = errors.New("movie not found")
	ErrMovieExists   = errors.New("movie exists")
)
