package rpc

import (
    "context"
    "fmt"
    "google.golang.org/grpc"
    "io"
    "net"
)

type GRPCServer struct {
    requestHandler func(request Request, reply *Reply) error
    grpcServer     *grpc.Server
    closer         io.Closer
}


//
//type GRPCSession struct {
//}
//
//func (s *GRPCServer) Send(stream GRPCService_SendServer) error {
//   // session := &GRPCSession{}
//    ctx := stream.Context()
//EXIT:
//    for {
//        select {
//            case <-ctx.Done():
//                break EXIT
//        default:
//            msg, err := stream.Recv()
//            if err != nil {
//                return err
//            }
//            fmt.Println(msg)
//        }
//    }
//    return nil
//}

func NewGRPCServer(handler func(request Request, reply *Reply) error) *GRPCServer {
    s := &GRPCServer{
        requestHandler: handler,
    }
    return s
}
func (s *GRPCServer) Call(ctx context.Context, in *GRPCRequest) (*GRPCResponse, error) {
    if s.requestHandler != nil {
        reply := &Reply{}
        err := s.requestHandler(Request{Action:in.Action, Data:in.Data}, reply)
        return &GRPCResponse{Code: int32(reply.Code), Data: reply.Data}, err
    }
    return &GRPCResponse{}, fmt.Errorf("handler not found")
}
func (s *GRPCServer) ListenServe(address string) error {
    ln, err := net.Listen("tcp", address)
    if err != nil {
        return err
    }

    grpcServer := grpc.NewServer()
    RegisterGRPCServiceServer(grpcServer, s)
    s.closer = ln
    s.grpcServer = grpcServer
    return grpcServer.Serve(ln)
}

func (s *GRPCServer) Close() error {
    return s.closer.Close()
}
