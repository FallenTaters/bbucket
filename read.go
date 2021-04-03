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
// GetterFunc receives a pointer to an object.
// Add the object to a slice defined in outside scope.
// Get your object of type T using: `*ptr.(*T)`
func (br Bucket) GetAll(dst interface{}, f GetterFunc) error {
	if f == nil {
		return ErrNilFuncPassed
	}

	return br.BucketView(func(b *bbolt.Bucket) error {
		return b.ForEach(func(_, v []byte) error {
			err := json.Unmarshal(v, dst)
			if err != nil {
				return err
			}

			f(dst)

			return nil
		})
	})
}
