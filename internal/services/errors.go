package services

import "errors"

var (
	ERROR_CONFLICT_TARGET = errors.New("The origin account and the target account are the same")
)
