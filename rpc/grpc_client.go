package rpc

import (
    "context"
    "google.golang.org/grpc"
    "time"
)

type GRPCClient struct {
    address string
    conn    *grpc.ClientConn
    cli     GRPCServiceClient
}

func NewGRPCClient(address string) *GRPCClient {
    cli := &GRPCClient{
        address: address,
    }
    return cli
}

func (c *GRPCClient) Connect() error {
    conn, err := grpc.Dial(c.address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*3))
    if err != nil {
        return err
    }
    c.conn = conn
    c.cli = NewGRPCServiceClient(conn)
    return nil
}

func (c *GRPCClient) Call(action string, data []byte) (Reply, error) {
    if c.cli == nil {
        return Reply{}, ErrorRPCClientInvalid
    }
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
    defer cancel()
    resp, err := c.cli.Call(ctx, &GRPCRequest{Action: action, Data: data})
    if err != nil {
        return Reply{}, err
    }
    return Reply{Code: int(resp.Code), Data: resp.Data}, nil
}
func (c *GRPCClient) OnMessage() {

}

func (c *GRPCClient) Disconnect() error {
    return c.conn.Close()
}
