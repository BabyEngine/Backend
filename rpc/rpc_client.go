package rpc

import (
    "github.com/BabyEngine/Backend/debugging"
    "net/rpc"
    "sync"
)

type Client struct {
    client  *rpc.Client
    mutex   sync.RWMutex
    status  int
    address string
}

const (
    rpcClientStatusInvalid    = -1
    rpcClientStatusConnecting = 1
    rpcClientStatusConnected  = 2
)

func NewClient(address string) (*Client, error) {
    client := &Client{
        client:  nil,
        status:  rpcClientStatusInvalid, // not connect
        address: address,
    }
    if err := client.connect(); err != nil {
        return nil, err
    }
    return client, nil
}

func (c *Client) Call(action string, data []byte) (r Reply, err error) {
    if c.client == nil || c.status != rpcClientStatusConnected {
        if err = c.connect(); err != nil {
            return
        }
    }
    if c.client == nil {
        err = ErrorRPCClientInvalid
        return
    }

    err = c.client.Call("RPC.Call", Request{Action: action, Data: data}, &r)
    if err != nil {
        c.status = rpcClientStatusInvalid
        if err == rpc.ErrShutdown {
            if err := c.connect(); err != nil {
                debugging.Log("rpc client auto connect error", err)
            }
        }
    }
    return
}

func (c *Client) Stop() error {
    err := c.client.Close()
    c.status = -1
    c.client = nil
    return err
}

func (c *Client) connect() error {
    if c.status == rpcClientStatusConnecting { // connecting
        return ErrorRPCClientConnecting
    }
    if c.status == rpcClientStatusConnected { // connected
        return nil
    }
    c.status = rpcClientStatusConnecting
    cli, err := rpc.DialHTTP("tcp", c.address)
    if err != nil {
        c.status = rpcClientStatusInvalid // not connect
        c.client = nil
        return err
    }
    c.status = rpcClientStatusConnected
    c.client = cli
    return nil
}
