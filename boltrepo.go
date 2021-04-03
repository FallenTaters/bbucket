package boltrepo

import "go.etcd.io/bbolt"

type BoltRepo struct {
	DB     *bbolt.DB
	Bucket []byte
}

func New(db *bbolt.DB, bucket []byte) BoltRepo {
	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		return err
	})
	if err != nil {
		panic(err)
	}

	return BoltRepo{
		DB:     db,
		Bucket: bucket,
	}
}

func (br BoltRepo) BucketView(f BucketFunc) error {
	return br.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(br.Bucket)
		if b == nil {
			return ErrBucketNotFound
		}

		return f(b)
	})
}

func (br BoltRepo) BucketUpdate(f BucketFunc) error {
	return br.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(br.Bucket)
		if b == nil {
			return ErrBucketNotFound
		}

		return f(b)
	})
}
