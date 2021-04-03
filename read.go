package boltrepo

import (
	"encoding/json"

	"go.etcd.io/bbolt"
)

func (br BoltRepo) Get(key []byte, dst interface{}) error {
	return br.BucketView(func(b *bbolt.Bucket) error {
		data := b.Get(key)
		if data == nil {
			return ErrObjectNotFound
		}

		return json.Unmarshal(data, dst)
	})
}

func (br BoltRepo) GetAll(dst interface{}, f GetterFunc) error {
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
