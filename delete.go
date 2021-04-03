package bbucket

import "go.etcd.io/bbolt"

// Delete deletes an object by key.
// If the key doesn't exist, it return ErrObjectNotFound
func (br Bucket) Delete(key []byte) error {
	return br.BucketUpdate(func(b *bbolt.Bucket) error {
		data := b.Get(key)
		if data == nil {
			return ErrObjectNotFound
		}

		return b.Delete(key)
	})
}
