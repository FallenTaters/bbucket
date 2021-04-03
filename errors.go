package bbucket

import (
	"errors"

	"go.etcd.io/bbolt"
)

var (
	ErrObjectNotFound      = errors.New("object not found")
	ErrObjectAlreadyExists = errors.New("object already exists")
	ErrKeyChanged          = errors.New("key change update not allowed")
	ErrBucketNotFound      = bbolt.ErrBucketNotFound
	ErrNilFuncPassed       = errors.New("nil function passed")
)
