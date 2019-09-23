package models

import "errors"

var (
	DB_BAD_RESPONSE = errors.New("The `db` return a response with errors")
)
