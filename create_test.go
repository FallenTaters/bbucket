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
	defer br.Close()

	t.Run(`marshal error`, func(t *testing.T) {
		assert := assert.New(t)

		err := br.Create(Itob(1), make(chan int))
		assert.Error(err)
	})

	t.Run(`successful creation`, func(t *testing.T) {
		assert := assert.New(t)

		expected := testStruct{
			ID:   "testCreate",
			Data: 123,
		}

		_, err := getTestStruct(br, expected.Key())
		assert.Eq(ErrObjectNotFound, err)

		err = br.Create(expected.Key(), expected)
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

		err = br.Create(expected.Key(), expected)
		assert.Eq(ErrObjectAlreadyExists, err)

		actual, err := getTestStruct(br, expected.Key())
		assert.NoError(err)

		assert.Eq(expected, actual)
	})
}

func TestNextSequence(t *testing.T) {
	assert := assert.New(t)
	br := getTestRepo()
	defer br.Close()

	i := br.NextSequence()
	assert.Eq(1, i)

	i = br.NextSequence()
	assert.Eq(2, i)
}
