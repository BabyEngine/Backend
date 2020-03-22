package kv

import (
    "github.com/BabyEngine/Backend/Debug"
    "github.com/boltdb/bolt"
)
type DB struct {
    db *bolt.DB
}
func OpenKVDB(path string) (*DB, error) {
    db, err := bolt.Open(path, 0600, nil)
    if err != nil {
        return nil, err
    }

    d := &DB{}
    d.db = db
    return d, nil
}

func (d *DB) Close()  {
    if d.db == nil {
        return
    }
    if err := d.db.Close(); err != nil {
        Debug.Log(err)
    }
    d.db = nil
}

func (d *DB) Update(bucketName string, key string, value []byte) error {
    err := d.db.Update(func(tx *bolt.Tx) error {
        bk, err := tx.CreateBucketIfNotExists([]byte(bucketName))
        if err != nil {
            return err
        }
        return bk.Put([]byte(key), value)
    })
    return err
}

func (d *DB) View(bucketName string, key string, cb func([]byte, error))  {
    var (
        value []byte
    )
    go func() {
        err := d.db.View(func(tx *bolt.Tx) error {
            bk := tx.Bucket([]byte(bucketName))
            if bk == nil {
                return bolt.ErrBucketNotFound
            }
            value = bk.Get([]byte(key))
            return nil
        })
        cb(value, err)
    }()
}

func (d *DB) Foreach(bucketName string, cb func(string, []byte))  {
    d.db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(bucketName))
        if b == nil {
            return bolt.ErrBucketNotFound
        }
        c := b.Cursor()
        for k, v := c.First(); k != nil; k, v = c.Next() {
            cb(string(k), v)
        }
        return nil
    })
}

func (d *DB) RemoveValue(bucketName string, key string)  {
    d.db.Update(func(tx *bolt.Tx) error {
        bk := tx.Bucket([]byte(bucketName))
        if bk == nil {
            return nil
        }
        bk.Delete([]byte(key))
        return nil
    })
}
func (d *DB) RemoveBucket(bucketName string)  {
    d.db.Update(func(tx *bolt.Tx) error {
        bk := tx.Bucket([]byte(bucketName))
        if bk == nil {
            return nil
        }
        bk.DeleteBucket([]byte(bucketName))
        return nil
    })
}