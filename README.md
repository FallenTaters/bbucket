# bbucket

## What does it do?

- Reduce boilerplate for go.etcd.io/bbolt for simple CRUD
- Allow for quick and dirty custom getters and filters

## What does it not do (yet)?

- Nested buckets
- Bulk update operations
- Bulk delete operations
- Bulk create operations
- Cursor-based use cases such as find()

# Get Started

```go
myBucketName := []byte("myBucket")

// open regular bbolt db
db, err := bbolt.Open(testPath, 0666, &bbolt.Options{Timeout: 1 * time.Second})
if err != nil {
    panic(err)
}

// make Bucket
myBucket := bbucket.New(db, myBucketName)
```

# Create

## Create

plain bbolt

```go
func create(obj Object) error {
    return db.Update(func(tx *bbolt.Tx) error {
        b := tx.Bucket(myBucketName)
        if b == nil {
            return errors.New("bucket not found")
        }

        if b.Get(obj.Key()) != nil {
            return errors.New("object already exists")
        }

        data, err := json.Marshal(obj)
        if err != nil {
            return err
        }

        return b.Put(obj.Key(), data)
    })
}
```

with bbucket

```go
func create(obj Object) error {
    return myBucket.Create(obj.Key(), obj)
}
```

## CreateAll

plain bbolt

```go
func createAll(objects []Object) error {
    return db.Update(func(tx *bbolt.Tx) error {
        b := tx.Bucket(myBucketName)
        if b == nil {
            return errors.New("bucket not found")
        }

        for _, obj := range objects {
            if b.Get(obj.Key()) != nil {
				return errors.New("object already exists")
			}

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
}
```

with bbucket

```go
func createAll(objects []Object) error {
	return br.CreateAll(objects, func(obj interface{}) (key []byte, err error) {
		return obj.(Object).Key(), nil
	})
}
```

# Get

## Get

plain bbolt

```go
func get(key []byte) (Object, error) {
    var obj Object
    return obj, db.View(func(tx *bbolt.Tx) error {
        b := tx.Bucket(myBucketName)
        if b == nil {
            return errors.New("bucket not found")
        }

        data := b.Get(key)
        if data == nil {
            return errors.New("object not found")
        }
        return json.Unmarshal(data, &obj)
    })
}
```

with bbucket

```go
func get(key []byte) (Object, error) {
    var obj Object
    return obj, myBucket.Get(key, &obj)
}
```

## Get All

plain bbolt

```go
func getByProp(prop int) ([]Object, error) {
    var objects []Object
    return objects, db.View(func(tx *bbolt.Tx) error {
        b := tx.Bucket(myBucketName)
        if b == nil {
            return errors.New("bucket not found")
        }

        return b.ForEach(func(k, v []byte) error {
            var obj Object
            err := json.Unmarshal(v, &obj)
            if err != nil {
                return err
            }

            if obj.prop == prop {
                objects = append(objects, obj)
            }

            return nil
        })
    })
}
```

with bbucket

```go
func getByProp(prop int) ([]Object, error) {
    var objects []Object
    return objects, myBucket.GetAll(&Object{}, func(obj interface{}) error {
        o := *obj.(*Object)

        if o.prop == prop {
            objects = append(objects, o)
        }

        return nil
    })
})
```

## Find

plain bbolt

```go
func findByProp(prop int) (Object, error) {
    var obj Object
    return obj, db.View(func(tx *bbolt.Tx) error {
        b := tx.Bucket(myBucketName)
        if b == nil {
            return errors.New("bucket not found")
        }

        c := b.Cursor()
        for _, v := c.First(); true; _, v = c.Next() {
            if v == nil {
                return errors.New("object not found")
            }

            err := json.Unmarshal(v, &obj)
            if err != nil {
                return err
            }

            if obj.prop == prop {
                return nil
            }
        }

    })
}
```

with bbucket

```go
func findByProp(prop int) (Object, error) {
    var obj Object
    return obj, myBucket.Find(&obj, func(_ []byte, ptr interface{}) (bool, err) {
        obj = *ptr.(*Object)
        return obj.prop == prop, nil
    })
}
```

# Update

## Update

plain bbolt

```go
func setProp(key []byte, value int) error {
    return db.Update(func(tx *bbolt.Tx) error {
        b := tx.Bucket(myBucketName)
        if b == nil {
            return errors.New("bucket not found")
        }

        data := b.Get(key)
        if data == nil {
            return errors.New("object not found")
        }

        var obj Object
        err := json.Unmarshal(data, &obj)
        if err != nil {
            return err
        }

        obj.prop = value

        if !bytes.Equal(obj.Key(), key) {
            err := b.Delete(key)
            if err != nil {
                return err
            }
        }

        data, err := json.Marshal(obj)
        if err != nil {
            return err
        }

        return b.Put(obj.Key(), data)
    })
}
```

with bbucket

```go
func setProp(key []byte, value int) error {
    return myBucket.Update(key, &Object{}, func(ptr interface{}) (interface{}, error) {
        obj := *ptr.(*Object)
        obj.prop = value
        return obj, nil
    })
}
```

## UpdateAll

```go
func setAllProp(objects []Object, value int) error {
	return br.UpdateAll(&Object{}, func(ptr interface{}) (key []byte, object interface{}, err error) {
        o := *obj.(*Object)
        o.prop = value
		return o.Key(), o, nil
	})
}
```

# Delete

plain bbolt

```go
func delete(key []byte) error {
    return db.Update(func(tx *bbolt.Tx) error {
        b := tx.Bucket(myBucketName)
        if b == nil {
            return errors.New("bucket not found")
        }

        data := b.Get(key)
        if data == nil {
            return errors.New("object not found")
        }

        return b.Delete(key)
    })
}
```

with bbucket

```go
func delete(key []byte) error {
    return myBucket.Delete(key)
}
```

# Coming Soon

## CreateAll([]Object) error

- create multiple objects from slice
- avoid multiple transactions for creating objects

## UpdateAll(dst interface{}, mutate MutateFunc)

- iterate over all objects in one transaction
- update where necessary

## DeleteAll(where WhereFunc)

- delete all records matched by your function in one transaction
