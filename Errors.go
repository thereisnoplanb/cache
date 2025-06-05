package cache

import "errors"

var ErrInvalidExpireAfter = errors.New("invalid expireAfter parameter")
var ErrKeyAlreadyExists = errors.New("key already exists")
var ErrKeyNotFound = errors.New("key not found")
