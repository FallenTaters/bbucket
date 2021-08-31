package bbucket

import (
	"errors"
	"fmt"

	"go.etcd.io/bbolt"
)

var (
	migrationBucket = []byte(`BBUCKET_MIGRATIONS`)

	ErrNotEnoughMigrations = errors.New(`not enough migrations`)
)

type Migration struct {
	Dst interface{}
	F   func(ptr interface{}) (key []byte, object interface{}, err error)
}

func (br Bucket) migrationState() int {
	var n int
	err := br.DB.Update(func(t *bbolt.Tx) error {
		b, err := t.CreateBucketIfNotExists(migrationBucket)
		if err != nil {
			return err
		}

		data := b.Get(br.Bucket)
		if data == nil {
			n = 0
			return nil
		}

		n = Btoi(data)
		return nil
	})
	if err != nil {
		panic(err)
	}

	return n
}

func (br Bucket) setMigrationState(n int) error {
	return br.DB.Update(func(t *bbolt.Tx) error {
		b, err := t.CreateBucketIfNotExists(migrationBucket)
		if err != nil {
			return err
		}
		return b.Put(br.Bucket, Itob(n))
	})
}

func (br Bucket) Migrate(migrations []Migration) error {
	currentState := br.migrationState()

	err := br.BucketUpdate(func(b *bbolt.Bucket) error {
		if len(migrations) < currentState {
			return fmt.Errorf(`%w: expected at least %d`, ErrNotEnoughMigrations, currentState)
		}

		for i := currentState; i < len(migrations); i++ {

			err := updateAll(b, migrations[i].Dst, migrations[i].F)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err == nil {
		return br.setMigrationState(len(migrations))
	}

	return err
}
