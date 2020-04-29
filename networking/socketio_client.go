package networking

import (
    "encoding/binary"
    "fmt"
    "github.com/BabyEngine/Backend/logger"
    socketio "github.com/googollee/go-socket.io"
    "sync/atomic"
    "time"
)

type mSocketIOClient struct {
    server     *mSocketIOServer
    id         int64
    conn       socketio.Conn
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
    msgChan    chan *mSocketIOClientMessage
}

type mSocketIOClientMessage struct {
    data []byte
    reply chan []byte
    evt string
}

func (c *mSocketIOClient) init() {
    c.stopChan = make(chan interface{}, 1)
    c.isStopRead = false
    c.info = ""
    c.msgChan = make(chan *mSocketIOClientMessage)
}

func (c *mSocketIOClient) Serve() {
    c.info = fmt.Sprintf("%v", c.conn.RemoteAddr())
    c.lastSeen = time.Now()
    defer func() {
        c.isStopRead = true
        c.status = "closed"
    }()
    if !c.opts.IsRawMode {
        // 发送一个 open 消息
        data := BuildMessage(OPCODE_OPEN, []byte(c.opts.Tag))
        //c.conn.WriteMessage(websocket.BinaryMessage, data)
        c.conn.Emit("message", data)
        c.status = "connected"
    }
    c.opts.Handler.OnNew(c)
EXITLOOP:
    for {
        select {
        case <-c.stopChan:
            break EXITLOOP
        default:
            msg := <-c.msgChan
            message := msg.data
            if c.opts.IsRawMode {
                n := uint64(len(message))
                atomic.AddUint64(&c.rp, 1)        // 收到数据包总数
                atomic.AddUint64(&c.rx, n)        // 收到的字节总数
                atomic.AddUint32(&c.qpsRCount, 1) // 收到的QPS
                c.server.onNetStat(0, 1, n)
                rs := c.opts.Handler.OnRequest(c, message) // socket io 不一样, 每一条消息必须有回复
                msg.reply <- rs
                continue
            }
            pkg, err := ParseMessage(message)
            if err != nil {
                msg.reply<-nil
                break EXITLOOP
            }
            // 统计
            atomic.AddUint64(&c.rp, 1)                   // 收到数据包总数
            atomic.AddUint64(&c.rx, uint64(pkg.DataLen)) // 收到的字节总数
            atomic.AddUint32(&c.qpsRCount, 1)            // 收到的QPS
            c.server.onNetStat(0, 1, uint64(pkg.DataLen))
            c.lastSeen = time.Now()
            switch pkg.OpCode {
            case OPCODE_OPEN:
            case OPCODE_CLOSE:
                msg.reply<-nil
                break EXITLOOP
            case OPCODE_PING:
                if pkg.DataLen == 4 {
                    val := binary.BigEndian.Uint32(pkg.Data)
                    c.latency = val
                }
                bin := BuildMessage(OPCODE_PONG, nil)
                c.conn.Emit(msg.evt, bin)
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
                    msg.reply <- bin
                    continue
                    //c.conn.Emit(TextMessage, bin)
                }
            case OPCODE_RESP:
            default:

            }
            logger.Debug("replay outside")
            msg.reply<-nil
        }
    }
}

func (c *mSocketIOClient) IsAlive() bool {
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

func (c *mSocketIOClient) Stop() {
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

func (c *mSocketIOClient) String() string {
    return fmt.Sprintf("%v %v", c.info, c.status)
}

func (c *mSocketIOClient) SendData(data []byte) error {
    return c.SendRawEvent("", OPCODE_DATA, data)
}

func (c *mSocketIOClient) SendRawEvent(e string, op OpCode, data []byte) error {
    if c.opts.IsRawMode {
        c.conn.Emit(e, data)
        return nil
    }
    bin := BuildMessage(op, data)
    n := len(bin)
    c.conn.Emit(e, data)
    atomic.AddUint64(&c.tp, 1)         // 发送的数据包总数
    atomic.AddUint64(&c.tx, uint64(n)) // 发送的数据包总字节数
    atomic.AddUint32(&c.qpsTCount, 1)  // 发送的QPS
    c.server.onNetStat(1, 1, uint64(n))
    return nil
}

func (c *mSocketIOClient) Close() {
    c.Stop()
}

func (c *mSocketIOClient) Id() int64 {
    return c.id
}
func (c *mSocketIOClient) SetId(id int64) {
    c.id = id
}
func (c *mSocketIOClient) RunCmd(cmd string, args []string) string {
    return ""
}
