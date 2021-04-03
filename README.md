# bbucket

## What does it do?

- Reduce boilerplate for go.etcd.io/bbolt for simple CRUD
- Allow for quick and dirty custom getters and filters

## What does it not do (yet)?

- Nested buckets
- Bulk update operations
- Bulk delete operations
- Bulk create operations
- Cursor-based use cases

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

# Examples

## simple getter by key

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

## get custom filtered slice

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
    return objects, myBucket.GetAll(&Object{}, func(obj interface{}) {
        o := *obj.(*Object)
        if o.prop == prop {
            objects = append(objects, o)
        }
    })
})
```

## update

## create

## delete
