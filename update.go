package bbucket

import (
	"bytes"
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
		if f == nil {
			return ErrNilFuncPassed
		}

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

func (br Bucket) UpdateAll(dst interface{}, f func(ptr interface{}) (key []byte, object interface{}, err error)) error {
	return br.BucketUpdate(func(b *bbolt.Bucket) error {
		return updateAll(b, dst, f)
	})
}

func updateAll(b *bbolt.Bucket, dst interface{}, f func(ptr interface{}) (key []byte, object interface{}, err error)) error {
	if f == nil {
		return ErrNilFuncPassed
	}

	type item struct {
		key  []byte
		data []byte
	}
	toBeDeleted := [][]byte{}
	toBePut := []item{}

	err := b.ForEach(func(originalKey, originalValue []byte) error {
		err := json.Unmarshal(originalValue, dst)
		if err != nil {
			return err
		}

		key, object, err := f(dst)
		if err != nil {
			return err
		}

		if key == nil { // delete item
			toBeDeleted = append(toBeDeleted, originalKey)
			return nil
		}

		data, err := json.Marshal(object)
		if err != nil {
			return err
		}

		if !bytes.Equal(originalKey, key) { // move item to new key, possibly with new value
			toBeDeleted = append(toBeDeleted, originalKey)
			toBePut = append(toBePut, item{key: key, data: data})
			return nil
		}

		if !bytes.Equal(originalValue, data) { // edit item but not key
			toBePut = append(toBePut, item{key, data})
		}

		return nil
	})
	if err != nil {
		return err
	}

	for _, key := range toBeDeleted {
		err := b.Delete(key)
		if err != nil {
			return err
		}
	}

	for _, item := range toBePut {
		err := b.Put(item.key, item.data)
		if err != nil {
			return err
		}
	}

	return nil
}
