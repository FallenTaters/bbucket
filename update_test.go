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

		err = br.Update(expected.Key(), &testStruct{}, func(objPtr interface{}) keyer {
			obj := *objPtr.(*testStruct)
			obj.Data = 9876
			return obj
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

		err = br.Update(expected.Key(), &testStruct{}, func(objPtr interface{}) keyer {
			return nil
		})

		assert.Eq(ErrObjectNotFound, err)
	})

	t.Run(`key change`, func(t *testing.T) {
		assert := assert.New(t)

		original := testStruct2
		actual, err := getTestStruct(br, original.Key())
		assert.NoError(err)
		assert.Eq(original, actual)

		expected := original
		expected.ID = "new_id"
		_, err = getTestStruct(br, expected.Key())
		assert.Eq(ErrObjectNotFound, err)

		err = br.Update(original.Key(), &testStruct{}, func(objPtr interface{}) keyer {
			obj := *objPtr.(*testStruct)
			obj.ID = "new_id"
			return obj
		})
		assert.NoError(err)

		actual, err = getTestStruct(br, expected.Key())
		assert.NoError(err)
		assert.Eq(expected, actual)

		_, err = getTestStruct(br, original.Key())
		assert.Eq(ErrObjectNotFound, err)
	})
}
