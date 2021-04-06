package bbucket

import "go.etcd.io/bbolt"

// Bucket has wrappers around Tx and Bucket to reduce boilerplate for simple CRUD operations
// It does not support sub-buckets
// A new Bucket should always be created using the New() constructor to ensure the bucket exists.
type Bucket struct {
	DB     *bbolt.DB
	Bucket []byte
}

// New returns a bbucket struct and ensures the bucket exists
// Panics for an invalid bucket name
func New(db *bbolt.DB, bucket []byte) Bucket {
	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		return err
	})
	if err != nil {
		panic(err)
	}

	return Bucket{
		DB:     db,
		Bucket: bucket,
	}
}

// BucketView is used internally and allows for custom implementations.
// It wraps DB.View() and Tx.Bucket()
func (br Bucket) BucketView(f func(*bbolt.Bucket) error) error {
	return br.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(br.Bucket)
		if b == nil {
			return ErrBucketNotFound
		}

		return f(b)
	})
}

// BucketUpdate is used internally and allows for custom implementations.
// It wraps DB.Update() and Tx.Bucket()
func (br Bucket) BucketUpdate(f func(*bbolt.Bucket) error) error {
	return br.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(br.Bucket)
		if b == nil {
			return ErrBucketNotFound
		}

		return f(b)
	})
}
