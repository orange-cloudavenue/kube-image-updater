package actions

import "errors"

var (
	// ErrActionNotFound is returned when an action is not found
	ErrActionNotFound = errors.New("action not found")

	// ErrEmptyNewTag is returned when the new tag is empty
	ErrEmptyNewTag = errors.New("new tag is empty")
)
