package constants

import "errors"

var (
	ErrCreate        = errors.New("failed to create")
	ErrUpdate        = errors.New("failed to update")
	ErrDelete        = errors.New("failed to delete")
	ErrNotFound      = errors.New("not found")
	ErrConnectOpenAI = errors.New("failed to connect openai")
	ErrInvalidUser   = errors.New("invalid user")
	ErrInvalidToken  = errors.New("token is invalid")
	ErrFailedLogin   = errors.New("failed to login")
)
