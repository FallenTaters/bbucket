package bbucket

import (
	"encoding/json"
	"reflect"

	"go.etcd.io/bbolt"
)

func (br Bucket) NextSequence() int {
	var i uint64
	err := br.BucketUpdate(func(b *bbolt.Bucket) error {
		j, err := b.NextSequence()
		i = j
		return err
	})
	if err != nil {
		panic(err)
	}

	return int(i)
}

// Create stores a new object in the bucket.
// If the key already exists, it returns ErrObjectAlreadyExists
func (br Bucket) Create(key []byte, obj interface{}) error {
	return br.BucketUpdate(func(b *bbolt.Bucket) error {
		if b.Get(key) != nil {
			return ErrObjectAlreadyExists
		}

		data, err := json.Marshal(obj)
		if err != nil {
			return err
		}

		return b.Put(key, data)
	})
}

// CreateAll stores multiple objects in the bucket.
// If any key already exists, it returns ErrObjectAlreadyExists
// First argument must be a slice of objects
// Second argument is a function that receives the object and should return the key
// Get your object by using obj.(Object)
// When relying on NextSequence for indexing, call it inside keyFunc
func (br Bucket) CreateAll(objs interface{}, keyFunc func(obj interface{}) (key []byte, err error)) error {
	return br.BucketUpdate(func(b *bbolt.Bucket) error {
		s := reflect.ValueOf(objs)
		if s.Kind() != reflect.Slice {
			return ErrNonSliceArgument
		}

		if s.IsNil() {
			return nil
		}

		for i := 0; i < s.Len(); i++ {
			obj := s.Index(i).Interface()

			key, err := keyFunc(obj)
			if err != nil {
				return err
			}

			if b.Get(key) != nil {
				return ErrObjectAlreadyExists
			}

			data, err := json.Marshal(obj)
			if err != nil {
				return err
			}

			err = b.Put(key, data)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
