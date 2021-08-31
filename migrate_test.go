package bbucket

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"git.fuyu.moe/Fuyu/assert"
	"go.etcd.io/bbolt"
)

type Type1 struct {
	K     int `json:"key"`
	Value int `json:"value"`
}

func (t Type1) Key() []byte {
	return Itob(t.K)
}

type Type2 struct {
	K     string `json:"new_key"`
	Value string `json:"new_value"`
}

func (t Type2) Key() []byte {
	return []byte(t.K)
}

func migrateTestData() []Type1 {
	return []Type1{
		{1, 1},
		{2, 202},
		{3, 300},
	}
}

func migrateExpected() []Type2 {
	return []Type2{
		{`1`, `1`},
		{`2`, `202`},
		{`3`, `300`},
	}
}

func testMigration() Migration {
	return Migration{
		Dst: &Type1{},
		F: func(ptr interface{}) (key []byte, object interface{}, err error) {
			o := *ptr.(*Type1)
			out := Type2{
				K:     fmt.Sprint(o.K),
				Value: fmt.Sprint(o.Value),
			}

			return out.Key(), out, nil
		},
	}
}

func TestMigrate(t *testing.T) {
	t.Run(`Successful Migration`, func(t *testing.T) {
		assert := assert.New(t)
		br := prepare()
		defer br.Close()

		err := br.BucketView(func(b *bbolt.Bucket) error {
			data := b.Get(migrationBucket)
			assert.Nil(data)
			return nil
		})
		assert.NoError(err)

		err = br.Migrate([]Migration{testMigration()})
		assert.NoError(err)

		actual := []Type2{}
		err = br.GetAll(&Type2{}, func(ptr interface{}) error {
			actual = append(actual, *ptr.(*Type2))
			return nil
		})
		assert.NoError(err)

		assert.Cmp(migrateExpected(), actual)
	})

	t.Run(`don't migrate if len(migrations) == migrationState`, func(t *testing.T) {
		assert := assert.New(t)
		br := prepare()
		defer br.Close()

		err := br.setMigrationState(1)
		assert.NoError(err)
		assert.Eq(1, br.migrationState())

		err = br.Migrate([]Migration{testMigration()})
		assert.NoError(err)

		actual := []Type1{}
		err = br.GetAll(&Type1{}, func(ptr interface{}) error {
			actual = append(actual, *ptr.(*Type1))
			return nil
		})
		assert.NoError(err)

		assert.Cmp(migrateTestData(), actual)
	})

	t.Run(`error during migration --> don't migrate, return same error`, func(t *testing.T) {
		assert := assert.New(t)
		br := prepare()
		defer br.Close()

		// custom error
		e := errors.New(`my_error`)
		err := br.Migrate([]Migration{{
			Dst: &Type1{},
			F: func(ptr interface{}) (key []byte, object interface{}, err error) {
				o := *ptr.(*Type1)
				out := Type2{
					K:     fmt.Sprint(o.K),
					Value: fmt.Sprint(o.Value),
				}

				return out.Key(), out, e
			},
		}})
		assert.True(errors.Is(err, e))
		assertNoMigration(assert, br)

		// marshal error
		err = br.Migrate([]Migration{{
			Dst: &Type1{},
			F: func(ptr interface{}) (key []byte, object interface{}, err error) {
				return []byte{1}, func() {}, nil
			},
		}})
		e = &json.UnsupportedTypeError{}
		assert.True(errors.As(err, &e))
		assertNoMigration(assert, br)

		// unmarshal error
		err = br.Migrate([]Migration{{
			Dst: Type1{},
			F: func(ptr interface{}) (key []byte, object interface{}, err error) {
				return []byte{1}, ptr.(Type1), nil
			},
		}})
		e = &json.UnsupportedTypeError{}
		assert.True(errors.As(err, &e))
		assertNoMigration(assert, br)
	})
}

func assertNoMigration(assert assert.Assert, br Bucket) {
	actual := []Type1{}
	err := br.GetAll(&Type1{}, func(ptr interface{}) error {
		actual = append(actual, *ptr.(*Type1))
		return nil
	})
	assert.NoError(err)

	assert.Cmp(migrateTestData(), actual)
}

func prepare() Bucket {
	db, err := bbolt.Open(testPath, 0666, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}

	_ = db.Update(func(tx *bbolt.Tx) error {
		_ = tx.DeleteBucket(testMigrateBucket)
		_ = tx.DeleteBucket(migrationBucket)
		_, _ = tx.CreateBucket(testMigrateBucket)
		return nil
	})

	br := Bucket{
		DB:     db,
		Bucket: testMigrateBucket,
	}

	err = br.BucketUpdate(func(b *bbolt.Bucket) error {
		for _, obj := range migrateTestData() {
			data, err := json.Marshal(obj)
			if err != nil {
				return err
			}
			err = b.Put(obj.Key(), data)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return br
}
