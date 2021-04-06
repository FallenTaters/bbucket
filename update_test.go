package bbucket

import (
	"testing"

	"git.fuyu.moe/Fuyu/assert"
)

func TestUpdate(t *testing.T) {
	br := getTestRepo()
	defer br.DB.Close()

	t.Run(`successful update`, func(t *testing.T) {
		assert := assert.New(t)

		expected := testStruct1
		actual, err := getTestStruct(br, expected.Key())
		assert.NoError(err)
		assert.Eq(expected, actual)

		expected.Data = 9876

		err = br.Update(expected.Key(), &testStruct{}, func(objPtr interface{}) (interface{}, error) {
			obj := *objPtr.(*testStruct)
			obj.Data = 9876
			return obj, nil
		})
		assert.NoError(err)

		actual, err = getTestStruct(br, expected.Key())
		assert.NoError(err)
		assert.Eq(expected, actual)
	})

	t.Run(`key not found returns ErrObjectNotFound`, func(t *testing.T) {
		assert := assert.New(t)

		expected := testStruct{
			ID: "blablabla",
		}

		_, err := getTestStruct(br, expected.Key())
		assert.Eq(ErrObjectNotFound, err)

		err = br.Update(expected.Key(), &testStruct{}, func(objPtr interface{}) (interface{}, error) {
			return nil, nil
		})

		assert.Eq(ErrObjectNotFound, err)
	})
}
