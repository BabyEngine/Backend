package networking

import (
    "encoding/binary"
    "fmt"
    "github.com/gorilla/websocket"
    "sync/atomic"
    "time"
)

type mWebsocketClient struct {
    server     *mWebsocketServer
    id         int64
    conn       *websocket.Conn
    opts       *Options
    qpsTCount  uint32 // transfer qps
    qpsRCount  uint32 // receive qps
    qpsT       uint32 // transfer qps
    qpsR       uint32 // receive qps
    latency    uint32 // latency
    tx         uint64 // transfer bytes
    rx         uint64 // receive bytes
    tp         uint64 // transfer packet
    rp         uint64 // receive packet
    stopChan   chan interface{}
    isStopRead bool
    isStop     bool
    lastSeen   time.Time // last seen time
    info       string
    status     string // connected, closed
}

func (c *mWebsocketClient) init() {
    c.stopChan = make(chan interface{}, 1)
    c.isStopRead = false
    c.info = ""
}

func (c *mWebsocketClient) Serve() {
    c.info = fmt.Sprintf("%v => %v", c.conn.RemoteAddr(), c.conn.LocalAddr())
    c.lastSeen = time.Now()
    defer func() {
        c.isStopRead = true
        c.status = "closed"
    }()
    if !c.opts.IsRawMode {
        // 发送一个 open 消息
        data := BuildMessage(OPCODE_OPEN, []byte(c.opts.Tag))
        c.conn.WriteMessage(websocket.BinaryMessage, data)
        c.status = "connected"
    }
    c.opts.Handler.OnNew(c)
EXITLOOP:
    for {
        select {
        case <-c.stopChan:
            break EXITLOOP
        default:
            _, message, err := c.conn.ReadMessage()
            //pkg, err := ReadMessage(c.conn)
            if err != nil {
                break EXITLOOP
            }
            if c.opts.IsRawMode {
                n := uint64(len(message))
                atomic.AddUint64(&c.rp, 1)        // 收到数据包总数
                atomic.AddUint64(&c.rx, n)        // 收到的字节总数
                atomic.AddUint32(&c.qpsRCount, 1) // 收到的QPS
                c.server.onNetStat(0, 1, n)
                c.opts.Handler.OnData(c, message)
                continue
            }
            pkg, err := ParseMessage(message)
            // 统计
            atomic.AddUint64(&c.rp, 1)                   // 收到数据包总数
            atomic.AddUint64(&c.rx, uint64(pkg.DataLen)) // 收到的字节总数
            atomic.AddUint32(&c.qpsRCount, 1)            // 收到的QPS
            c.server.onNetStat(0, 1, uint64(pkg.DataLen))
            c.lastSeen = time.Now()
            switch pkg.OpCode {
            case OPCODE_OPEN:
            case OPCODE_CLOSE:
                break EXITLOOP
            case OPCODE_PING:
                if pkg.DataLen == 4 {
                    val := binary.BigEndian.Uint32(pkg.Data)
                    c.latency = val
                }
                bin := BuildMessage(OPCODE_PONG, nil)
                if err := c.conn.WriteMessage(websocket.BinaryMessage, bin); err != nil {
                    c.opts.Handler.OnError(c, err)
                    break EXITLOOP
                }
            case OPCODE_PONG:
            case OPCODE_DATA:
                c.opts.Handler.OnData(c, pkg.Data)
            case OPCODE_TURN:
            case OPCODE_NOOP:
                c.Stop()
            case OPCODE_REQ:
                if pkg.DataLen > 4 {
                    reqId := make([]byte, 4)
                    copy(reqId, pkg.Data[:4])
                    body := pkg.Data[4:]
                    respData := c.opts.Handler.OnRequest(c, body)
                    if respData == nil {
                        respData = []byte{}
                    }
                    respFullData := append(reqId, respData...)
                    bin := BuildMessage(OPCODE_RESP, respFullData)
                    if err := c.conn.WriteMessage(websocket.BinaryMessage, bin); err != nil {
                        c.opts.Handler.OnError(c, err)
                        break EXITLOOP
                    }
                }
            case OPCODE_RESP:
            default:

            }
        }
    }
}

func (c *mWebsocketClient) IsAlive() bool {
    if c.isStopRead {
        return false
    }
    if c.opts.IsRawMode {
        return true
    }
    if time.Now().Sub(c.lastSeen) > c.opts.TTL {
        c.status = "timeout"
        c.Stop()
        return false
    }
    return true
}

func (c *mWebsocketClient) Stop() {
    if c.isStopRead {
        return
    }
    if c.isStop {
        return
    }
    c.isStop = true
    _ = c.conn.Close()
    go func() {
    EXITLOOP:
        for {
            timeout := time.After(time.Second)
            select {
            case c.stopChan <- 1:
                // notify ok
                break EXITLOOP
            case <-timeout:
                // notify failed
                if c.isStopRead { // but somewhere exit read loop
                    // exit ok
                    break EXITLOOP
                }
                // continue notify exit
            }
        }
    }()
}

func (c *mWebsocketClient) String() string {
    return fmt.Sprintf("%v %v", c.info, c.status)
}

func (c *mWebsocketClient) SendData(data []byte) error {
    return c.SendRaw(OPCODE_DATA, data)

}
func (c *mWebsocketClient) SendRaw(op OpCode, data []byte) error {
    if c.opts.IsRawMode {
        return c.conn.WriteMessage(websocket.BinaryMessage, data)
    }
    bin := BuildMessage(op, data)
    n := len(bin)
    if err := c.conn.WriteMessage(websocket.BinaryMessage, bin); err != nil {
        atomic.AddUint64(&c.tp, 1)         // 发送的数据包总数
        atomic.AddUint64(&c.tx, uint64(n)) // 发送的数据包总字节数
        atomic.AddUint32(&c.qpsTCount, 1)  // 发送的QPS
        c.server.onNetStat(1, 1, uint64(n))
        return nil
    } else {
        return err
    }
}
func (c *mWebsocketClient) Close() {
    c.Stop()
}

func (c *mWebsocketClient) Id() int64 {
    return c.id
}
func (c *mWebsocketClient) SetId(id int64)  {
    c.id = id
}
func (c *mWebsocketClient) RunCmd(cmd string, args[] string) string {
    return ""
}
