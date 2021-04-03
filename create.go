package boltrepo

import (
	"encoding/json"

	"go.etcd.io/bbolt"
)

func (br BoltRepo) Create(obj keyer) error {
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
