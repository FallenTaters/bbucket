package bbucket

import (
	"errors"
	"strings"
	"testing"

	"git.fuyu.moe/Fuyu/assert"
	"github.com/google/go-cmp/cmp/cmpopts"
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

func TestUpdateAll(t *testing.T) {
	t.Run(`nil func passed`, func(t *testing.T) {
		assert := assert.New(t)
		br := getTestRepo()
		defer br.Close()

		err := br.UpdateAll(&testStruct{}, nil)

		assert.Eq(ErrNilFuncPassed, err)
	})

	t.Run(`don't save if any error occurs`, func(t *testing.T) {
		assert := assert.New(t)
		br := getTestRepo()
		defer br.Close()

		// unmarshal error
		err := br.UpdateAll(testStruct{}, func(ptr interface{}) (key []byte, object interface{}, err error) {
			return nil, nil, nil
		})
		assert.Error(err)
		assertUnchanged(assert, br)

		// marshal error
		err = br.UpdateAll(&testStruct{}, func(ptr interface{}) (key []byte, object interface{}, err error) {
			return []byte(`bla`), func() {}, nil
		})
		assert.Error(err)
		assertUnchanged(assert, br)

		// custom error in function
		err = br.UpdateAll(&testStruct{}, func(ptr interface{}) (key []byte, object interface{}, err error) {
			return []byte(`bla`), testStruct{}, errors.New(`sdafkjaslkdfj`)
		})
		assert.Error(err)
		assertUnchanged(assert, br)
	})

	t.Run(`delete update`, func(t *testing.T) {
		assert := assert.New(t)
		br := getTestRepo()
		defer br.Close()

		o := testStruct{}
		err := br.Get([]byte(`ABC`), &o)
		assert.NoError(err)

		err = br.UpdateAll(&testStruct{}, func(ptr interface{}) (key []byte, object interface{}, err error) {
			o := *ptr.(*testStruct)
			if o.ID == `ABC` {
				return nil, nil, nil
			}
			return o.Key(), o, nil
		})
		assert.NoError(err)

		err = br.Get([]byte(`ABC`), &o)
		assert.Eq(ErrObjectNotFound, err)
	})

	t.Run(`update key`, func(t *testing.T) {
		assert := assert.New(t)
		br := getTestRepo()
		defer br.Close()

		o := testStruct{}
		err := br.Get([]byte(`bla`), &o)
		assert.Eq(ErrObjectNotFound, err)

		err = br.UpdateAll(&testStruct{}, func(ptr interface{}) (key []byte, object interface{}, err error) {
			o := *ptr.(*testStruct)
			if o.ID == `ABC` {
				return []byte(`bla`), o, nil
			}
			return o.Key(), o, nil
		})
		assert.NoError(err)

		err = br.Get([]byte(`bla`), &o)
		assert.NoError(err)
	})

	t.Run(`update value`, func(t *testing.T) {
		assert := assert.New(t)
		br := getTestRepo()
		defer br.Close()

		err := br.UpdateAll(&testStruct{}, func(ptr interface{}) (key []byte, object interface{}, err error) {
			o := *ptr.(*testStruct)
			if o.ID == `ABC` {
				o.Data = 99999
			}
			return o.Key(), o, nil
		})
		assert.NoError(err)

		o := testStruct{}
		err = br.Get([]byte(`ABC`), &o)
		assert.NoError(err)
		assert.Eq(99999, o.Data)
	})

	t.Run(`update key and value`, func(t *testing.T) {
		assert := assert.New(t)
		br := getTestRepo()
		defer br.Close()

		o := testStruct{}
		err := br.Get([]byte(`bla`), &o)
		assert.Eq(ErrObjectNotFound, err)

		err = br.UpdateAll(&testStruct{}, func(ptr interface{}) (key []byte, object interface{}, err error) {
			o := *ptr.(*testStruct)
			if o.ID == `ABC` {
				o.Data = 99999
				return []byte(`bla`), o, nil
			}
			return o.Key(), o, nil
		})
		assert.NoError(err)

		err = br.Get([]byte(`bla`), &o)
		assert.Eq(99999, o.Data)
		assert.NoError(err)
	})
}

func assertUnchanged(assert assert.Assert, br Bucket) {
	have := getAllTestStructs(br)
	assert.Cmp(testData, have, cmpopts.SortSlices(func(a, b testStruct) bool {
		return strings.Compare(a.ID, b.ID) < 0
	}))
}
