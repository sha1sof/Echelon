package storage

import "errors"

var (
	ErrVideoIDExist    = errors.New("videoID already exist")
	ErrVideoIDNotFound = errors.New("videoID not found")
)
