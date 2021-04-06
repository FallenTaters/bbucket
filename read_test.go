package bbucket

import (
	"testing"

	"git.fuyu.moe/Fuyu/assert"
)

func TestGet(t *testing.T) {
	br := getTestRepo()
	defer br.DB.Close()

	t.Run(`found`, func(t *testing.T) {
		assert := assert.New(t)

		expected := testStruct1

		var actual testStruct
		err := br.Get(expected.Key(), &actual)

		assert.NoError(err)
		assert.Cmp(expected, actual)
	})

	t.Run(`not found`, func(t *testing.T) {
		assert := assert.New(t)

		var actual testStruct
		err := br.Get([]byte("blablabla"), &actual)

		assert.Eq(ErrObjectNotFound, err)
	})
}

func TestGetAll(t *testing.T) {
	br := getTestRepo()
	defer br.DB.Close()

	t.Run(`get all`, func(t *testing.T) {
		assert := assert.New(t)

		expected := testData
		var actual []testStruct
		err := br.GetAll(&testStruct{}, func(obj interface{}) error {
			actual = append(actual, *obj.(*testStruct))
			return nil
		})
		assert.NoError(err)

		assert.Cmp(expected, actual)
	})
}
