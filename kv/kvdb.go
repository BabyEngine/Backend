package kv

import (
    "github.com/BabyEngine/Backend/logger"
    "github.com/boltdb/bolt"
)
type KVDB struct {
    db *bolt.DB
}
func OpenKVDB(path string) (*KVDB, error) {
    db, err := bolt.Open(path, 0600, nil)
    if err != nil {
        return nil, err
    }

    d := &KVDB{}
    d.db = db
    return d, nil
}

func (d *KVDB) Close()  {
    if d.db == nil {
        return
    }
    if err := d.db.Close(); err != nil {
        logger.Debug(err)
    }
    d.db = nil
}

func (d *KVDB) Update(bucketName string, key string, value []byte) error {
    err := d.db.Update(func(tx *bolt.Tx) error {
        bk, err := tx.CreateBucketIfNotExists([]byte(bucketName))
        if err != nil {
            return err
        }
        return bk.Put([]byte(key), value)
    })
    return err
}

func (d *KVDB) View(bucketName string, key string, cb func([]byte, error))  {
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

func (d *KVDB) Foreach(bucketName string, cb func(string, []byte))  {
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

func (d *KVDB) RemoveValue(bucketName string, key string)  {
    d.db.Update(func(tx *bolt.Tx) error {
        bk := tx.Bucket([]byte(bucketName))
        if bk == nil {
            return nil
        }
        bk.Delete([]byte(key))
        return nil
    })
}
func (d *KVDB) RemoveBucket(bucketName string)  {
    d.db.Update(func(tx *bolt.Tx) error {
        bk := tx.Bucket([]byte(bucketName))
        if bk == nil {
            return nil
        }
        bk.DeleteBucket([]byte(bucketName))
        return nil
    })
}