package bbucket

import (
	"errors"
	"testing"

	"git.fuyu.moe/Fuyu/assert"
)

func TestUpdate(t *testing.T) {
	br := getTestRepo()
	defer br.Close()

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

	t.Run(`invalid destination`, func(t *testing.T) {
		assert := assert.New(t)

		err := br.Get(testStruct1.Key(), &testStruct{})
		assert.NoError(err)

		err = br.Update(testStruct1.Key(), make(chan int), func(obj interface{}) (interface{}, error) {
			return nil, nil
		})
		assert.Error(err)
	})

	t.Run(`nil function`, func(t *testing.T) {
		assert := assert.New(t)

		err := br.Get(testStruct1.Key(), &testStruct{})
		assert.NoError(err)

		err = br.Update(testStruct1.Key(), &testStruct{}, nil)
		assert.Eq(ErrNilFuncPassed, err)
	})

	t.Run(`pass error through`, func(t *testing.T) {
		assert := assert.New(t)

		myErr := errors.New(`custom error`)

		err := br.Update(testStruct1.Key(), &testStruct{}, func(_ interface{}) (interface{}, error) {
			return nil, myErr
		})
		assert.Eq(myErr, err)
	})

	t.Run(`marshal error`, func(t *testing.T) {
		assert := assert.New(t)

		err := br.Get(testStruct1.Key(), &testStruct{})
		assert.NoError(err)

		err = br.Update(testStruct1.Key(), &testStruct{}, func(_ interface{}) (interface{}, error) {
			return make(chan int), nil
		})
		assert.Error(err)
	})
}
