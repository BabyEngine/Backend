package rpc

import (
    "fmt"
    "net"
    "net/rpc"
    "net/rpc/jsonrpc"
    "sync"
    "sync/atomic"
    "time"
)

type JSONRPCClient struct {
    mutex   sync.RWMutex
    address string
    client *rpc.Client
    status int32
}

func NewJSONRPCClient(address string) *JSONRPCClient {
    client := &JSONRPCClient{
        address: address,
    }
    return client
}

func (c *JSONRPCClient) Connect() error {
    atomic.StoreInt32(&c.status, 1)
    conn, err := net.DialTimeout("tcp", c.address, time.Second * 30)
    if err != nil {
        atomic.StoreInt32(&c.status, 3)
        return err
    }
    c.client = jsonrpc.NewClient(conn)
    atomic.StoreInt32(&c.status, 2)
    return nil
}
func (c *JSONRPCClient) reconnect()  {
    c.Connect()
}
func (c *JSONRPCClient) Disconnect() error {
    return nil
}

func (c *JSONRPCClient) Call(action string, data []byte) (r Reply, err error) {
    req := &Request{Action:action, Data:data}
    rsp := new(Reply)
    err = c.client.Call("JsonRpcServer.Invoke", req, &rsp)
    if err != nil {
        b := atomic.CompareAndSwapInt32(&c.status, 2, 1)
        if b {
            fmt.Println("reconnect...")
            go c.reconnect()
        } else {
            switch atomic.LoadInt32(&c.status) {
            case 0: //  init
            case 1: // connecting
            case 2: // connected
            case 3: // error
                go c.reconnect()
            }
        }

    }
    r.Data = rsp.Data
    r.Code = rsp.Code
    return
}

func (c *JSONRPCClient) Stop() error {
    c.client.Close()
    return nil
}
