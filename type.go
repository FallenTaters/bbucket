package bbucket

import "go.etcd.io/bbolt"

type Keyer interface {
	Key() []byte
}

// BucketFunc receives a bucket to operate on.
// It is intended for custom implementations not covered by bbucket.
type BucketFunc func(*bbolt.Bucket) error

// GetterFunc receives a pointer to an object.
// Add the object to a slice defined in outside scope.
// Get your object of type T using: `*ptr.(*T)`
type GetterFunc func(ptr interface{})

// MutateFunc receives a pointer to an object to be modified.
// It should return the modified object, not a pointer.
// Get your object of type T using: `*ptr.(*T)`
type MutateFunc func(ptr interface{}) Keyer
