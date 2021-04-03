package bbucket

import (
	"bytes"
	"encoding/json"

	"go.etcd.io/bbolt"
)

// Update saves changes to an object made in a function.
// If the key does not exist, it return ErrObjectNotFound
// MutateFunc receives a pointer to an object to be modified.
// It should return the modified object, not a pointer.
// Get your object of type T using: `*ptr.(*T)`
func (br Bucket) Update(key []byte, dst interface{}, mutate MutateFunc) error {
	return br.BucketUpdate(func(b *bbolt.Bucket) error {
		data := b.Get(key)
		if data == nil {
			return ErrObjectNotFound
		}

		err := json.Unmarshal(data, dst)
		if err != nil {
			return err
		}

		newObj := mutate(dst)
		newKey := newObj.(keyer).Key()

		if !bytes.Equal(newKey, key) {
			err = b.Delete(key)
			if err != nil {
				return err
			}
		}

		data, err = json.Marshal(newObj)
		if err != nil {
			return err
		}

		return b.Put(newKey, data)
	})
}
