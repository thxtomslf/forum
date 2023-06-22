package errors

import "errors"

var (
	ErrDuplicate     = errors.New("duplicate")
	ErrConflict      = errors.New("user conflict")
	ErrUserNotFound  = errors.New("user not found")
	ErrForumNotFound = errors.New("forum not found")
)
