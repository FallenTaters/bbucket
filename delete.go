package boltrepo

import "go.etcd.io/bbolt"

func (br BoltRepo) Delete(key []byte) error {
	return br.BucketUpdate(func(b *bbolt.Bucket) error {
		data := b.Get(key)
		if data == nil {
			return ErrObjectNotFound
		}

		return b.Delete(key)
	})
}
