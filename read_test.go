package bbucket

import (
	"errors"
	"testing"

	"git.fuyu.moe/Fuyu/assert"
)

func TestGet(t *testing.T) {
	br := getTestRepo()
	defer br.Close()

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
	defer br.Close()

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

	t.Run(`nil function`, func(t *testing.T) {
		assert := assert.New(t)

		err := br.GetAll(&testStruct{}, nil)
		assert.Eq(ErrNilFuncPassed, err)
	})

	t.Run(`invalid destination`, func(t *testing.T) {
		assert := assert.New(t)

		c := make(chan int)

		err := br.GetAll(c, func(obj interface{}) error {
			return nil
		})
		assert.Error(err)
	})
}

func TestFind(t *testing.T) {
	br := getTestRepo()
	defer br.Close()

	t.Run(`find returning found = true`, func(t *testing.T) {
		assert := assert.New(t)

		var obj testStruct
		err := br.Find(&testStruct{}, func(key []byte, ptr interface{}) (found bool, err error) {
			obj = *ptr.(*testStruct)
			return obj == testStruct2, nil
		})
		assert.NoError(err)
		assert.Eq(obj, testStruct2)
	})

	t.Run(`find returns ErrObjectNotFound`, func(t *testing.T) {
		assert := assert.New(t)

		var obj testStruct
		err := br.Find(&testStruct{}, func(key []byte, ptr interface{}) (found bool, err error) {
			obj = *ptr.(*testStruct)
			return false, nil
		})
		assert.Eq(ErrObjectNotFound, err)
		assert.Eq(obj, testStruct3)
	})

	t.Run(`error interrupts and is passed through`, func(t *testing.T) {
		assert := assert.New(t)

		var obj testStruct
		myErr := errors.New("custom error")
		err := br.Find(&testStruct{}, func(key []byte, ptr interface{}) (found bool, err error) {
			obj = *ptr.(*testStruct)
			return false, myErr
		})
		assert.Eq(myErr, err)
		assert.Eq(obj, testStruct1)
	})

	t.Run(`nil function`, func(t *testing.T) {
		assert := assert.New(t)

		err := br.Find(&testStruct{}, nil)
		assert.Eq(ErrNilFuncPassed, err)
	})

	t.Run(`invalid destination`, func(t *testing.T) {
		assert := assert.New(t)

		c := make(chan int)

		err := br.Find(c, func(key []byte, ptr interface{}) (found bool, err error) {
			return false, nil
		})
		assert.Error(err)
	})
}
