package bbucket

import (
	"encoding/json"

	"go.etcd.io/bbolt"
)

// Create stores a new object in the bucket.
// If the key already exists, it returns ErrObjectAlreadyExists
func (br Bucket) Create(obj Keyer) error {
	return br.BucketUpdate(func(b *bbolt.Bucket) error {
		data := b.Get(obj.Key())
		if data != nil {
			return ErrObjectAlreadyExists
		}

		data, err := json.Marshal(obj)
		if err != nil {
			return err
		}

		return b.Put(obj.Key(), data)
	})
}
