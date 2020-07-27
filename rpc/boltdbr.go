package rpc

import (
    "github.com/DGHeroin/boltdbr"
    "net"
    "net/rpc"
    "net/rpc/jsonrpc"
    "time"
)

type BoltDBRClient struct {
    conn    net.Conn
    address string
    client  *rpc.Client
    token   string
}

func NewBoltDBR(token string) *BoltDBRClient {
    c := &BoltDBRClient{
        token: token,
    }
    return c
}
func (c *BoltDBRClient) Connect(address string, cb func(error)) {
    c.address = address
    go func() {
        conn, err := net.DialTimeout("tcp", address, time.Second*30)
        if err != nil {
            cb(err)
            return
        }
        c.client = jsonrpc.NewClient(conn)
        cb(nil)
    }()
}

func (c *BoltDBRClient) reconnect() {
    c.Connect(c.address, func(err error) {

    })
}

func (c *BoltDBRClient) Get(bucket string, key string, cb func(string, string)) {
    go func() {
        q := &boltdbr.Query{Bucket: []byte(bucket), Token: c.token, Key: []byte(key)}
        r := new(boltdbr.Response)
        err := c.client.Call("BoltDBR.Get", q, &r)
        if err != nil {
            cb("", err.Error())
            c.reconnect()
            return
        }
        cb(string(r.Value), r.Error)
    }()
}

func (c *BoltDBRClient) Set(bucket string, key string, value string, cb func(string)) {
    go func() {
        q := &boltdbr.Query{Bucket: []byte(bucket), Token: c.token, Key: []byte(key), Value: []byte(value)}
        r := new(boltdbr.Response)
        err := c.client.Call("BoltDBR.Set", q, &r)
        if err != nil {
            cb(err.Error())
            c.reconnect()
            return
        }
        cb(r.Error)
    }()
}

func (c *BoltDBRClient) Delete(bucket string, key string, cb func(string)) {
    go func() {
        q := &boltdbr.Query{Bucket: []byte(bucket), Token: c.token, Key: []byte(key)}
        r := new(boltdbr.Response)
        err := c.client.Call("BoltDBR.Delete", q, &r)
        if err != nil {
            cb(err.Error())
            c.reconnect()
            return
        }
        cb(r.Error)
    }()
}
