package rpc

import (
    "bytes"
    "fmt"
    "github.com/gorilla/rpc/json"
    "net/http"
    "sync"
)

type JSONRPCClient struct {
    mutex   sync.RWMutex
    address string
}

func NewJSONRPCClient(address string) (*JSONRPCClient, error) {
    client := &JSONRPCClient{
        address: address,
    }
    return client, nil
}

func (c *JSONRPCClient) Call(action string, data []byte) (r Reply, err error) {
    var(
        message []byte
        resp *http.Response
    )
    message, err = json.EncodeClientRequest("RPC.Call", Request{Action: action, Data: data})
    if err != nil {
        return
    }
    url := fmt.Sprintf("%s/jsonrpc", c.address)
    resp, err = http.Post(url, "application/json", bytes.NewReader(message))
    if err != nil {
        return
    }
    err = json.DecodeClientResponse(resp.Body, &r)
    return
}

func (c *JSONRPCClient) Stop() error {
    return nil
}
