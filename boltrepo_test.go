package boltrepo

import (
	"encoding/json"
	"time"

	"go.etcd.io/bbolt"
)

var (
	testPath   = "test.db"
	testBucket = []byte("test")
)

type testStruct struct {
	ID   string `json:"a"`
	Data int    `json:"b"`
}

func (ts testStruct) Key() []byte {
	return []byte(ts.ID)
}

var (
	testStruct1 = testStruct{
		"ABC",
		123,
	}
	testStruct2 = testStruct{
		"BCD",
		234,
	}
	testStruct3 = testStruct{
		"CDE",
		345,
	}
)

var testData = []testStruct{
	testStruct1, testStruct2, testStruct3,
}

func getTestRepo() BoltRepo {
	db, err := bbolt.Open(testPath, 0666, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}

	_ = db.Update(func(tx *bbolt.Tx) error {
		_ = tx.DeleteBucket(testBucket)
		_, _ = tx.CreateBucket(testBucket)
		return nil
	})

	br := BoltRepo{
		DB:     db,
		Bucket: []byte("test"),
	}

	err = br.BucketUpdate(func(b *bbolt.Bucket) error {
		for _, obj := range testData {
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
