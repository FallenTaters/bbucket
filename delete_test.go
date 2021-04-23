package bbucket

import (
	"testing"

	"git.fuyu.moe/Fuyu/assert"
)

func TestDelete(t *testing.T) {
	br := getTestRepo()
	defer br.Close()

	t.Run(`successful deletion`, func(t *testing.T) {
		assert := assert.New(t)

		expected := testStruct1
		_, err := getTestStruct(br, expected.Key())
		assert.NoError(err)

		err = br.Delete(expected.Key())
		assert.NoError(err)

		_, err = getTestStruct(br, expected.Key())
		assert.Eq(ErrObjectNotFound, err)
	})

	t.Run(`object not found`, func(t *testing.T) {
		assert := assert.New(t)

		expected := testStruct4
		_, err := getTestStruct(br, expected.Key())
		assert.Eq(ErrObjectNotFound, err)

		err = br.Delete(expected.Key())
		assert.Eq(ErrObjectNotFound, err)
	})
}
