package bbucket

import (
	"encoding/json"
	"testing"

	"git.fuyu.moe/Fuyu/assert"
	"go.etcd.io/bbolt"
)

func getTestStruct(br Bucket, key []byte) (testStruct, error) {
	var t testStruct

	return t, br.BucketView(func(b *bbolt.Bucket) error {
		data := b.Get(key)
		if data == nil {
			return ErrObjectNotFound
		}

		return json.Unmarshal(data, &t)
	})
}

func TestCreate(t *testing.T) {
	br := getTestRepo()
	defer br.DB.Close()

	t.Run(`successful creation`, func(t *testing.T) {
		assert := assert.New(t)

		expected := testStruct{
			ID:   "testCreate",
			Data: 123,
		}

		_, err := getTestStruct(br, expected.Key())
		assert.Eq(ErrObjectNotFound, err)

		err = br.Create(expected)
		assert.NoError(err)

		actual, err := getTestStruct(br, expected.Key())
		assert.NoError(err)

		assert.Eq(expected, actual)
	})

	t.Run(`attempt duplicate key creation`, func(t *testing.T) {
		assert := assert.New(t)

		expected := testStruct{
			ID:   "testCreate",
			Data: 123,
		}

		_, err := getTestStruct(br, expected.Key())
		assert.NoError(err)

		err = br.Create(expected)
		assert.Eq(ErrObjectAlreadyExists, err)

		actual, err := getTestStruct(br, expected.Key())
		assert.NoError(err)

		assert.Eq(expected, actual)
	})
}
