package bbucket

import (
	"encoding/json"

	"go.etcd.io/bbolt"
)

// Get scans a single object by key
// If the key is unknown, it returns ErrObjectNotFound
func (br Bucket) Get(key []byte, dst interface{}) error {
	return br.BucketView(func(b *bbolt.Bucket) error {
		data := b.Get(key)
		if data == nil {
			return ErrObjectNotFound
		}

		return json.Unmarshal(data, dst)
	})
}

// GetAll iterates over all objects in the bucket.
// f receives a pointer to an object.
// Add the object to a slice defined in outside scope.
// Get your object of type T using: `*ptr.(*T)`
func (br Bucket) GetAll(dst interface{}, f func(ptr interface{}) error) error {
	if f == nil {
		return ErrNilFuncPassed
	}

	return br.BucketView(func(b *bbolt.Bucket) error {
		return b.ForEach(func(_, v []byte) error {
			err := json.Unmarshal(v, dst)
			if err != nil {
				return err
			}

			return f(dst)
		})
	})
}

// Find uses a cursor to iterate over all objects in the bucket.
// It stops when f returns found = true or a non-nil error.
// If it reaches the end, it returns ErrObjectNotFound
// Get your object of type T using: `*ptr.(*T)`
func (br Bucket) Find(dst interface{}, f func(key []byte, ptr interface{}) (found bool, err error)) error {
	if f == nil {
		return ErrNilFuncPassed
	}

	return br.BucketView(func(b *bbolt.Bucket) error {
		c, found := b.Cursor(), false
		for k, v := c.First(); !found; k, v = c.Next() {
			if k == nil {
				return ErrObjectNotFound
			}

			err := json.Unmarshal(v, dst)
			if err != nil {
				return err
			}

			found, err = f(k, dst)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
