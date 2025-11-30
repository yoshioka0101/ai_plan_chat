package entity

import "errors"

var (
	// Item errors
	ErrEmptyUserID        = errors.New("user ID cannot be empty")
	ErrEmptyTitle         = errors.New("title cannot be empty")
	ErrEmptyOriginalInput = errors.New("original input cannot be empty")
	ErrEmptyAIModel       = errors.New("AI model cannot be empty")
	ErrInvalidItemType    = errors.New("invalid item type")
	ErrItemNotFound       = errors.New("item not found")
)
