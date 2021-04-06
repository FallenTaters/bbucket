package bbucket

import (
	"encoding/json"

	"go.etcd.io/bbolt"
)

// Update saves changes to an object made in a function.
// If the key does not exist, it return ErrObjectNotFound
// f receives a pointer to an object to be modified.
// It should return the modified object, not a pointer.
// Get your object of type T using: `*ptr.(*T)`
func (br Bucket) Update(key []byte, dst interface{}, f func(ptr interface{}) (object interface{}, err error)) error {
	return br.BucketUpdate(func(b *bbolt.Bucket) error {
		data := b.Get(key)
		if data == nil {
			return ErrObjectNotFound
		}

		err := json.Unmarshal(data, dst)
		if err != nil {
			return err
		}

		obj, err := f(dst)
		if err != nil {
			return err
		}

		data, err = json.Marshal(obj)
		if err != nil {
			return err
		}

		return b.Put(key, data)
	})
}
