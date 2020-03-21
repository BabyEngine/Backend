package networking

import (
    "bytes"
    "encoding/binary"
    "io"
    "net"
    "time"
)

type OpCode uint8

const (
    OPCODE_OPEN  = OpCode(0) // 链路打开
    OPCODE_CLOSE = OpCode(1) // 链路关闭
    OPCODE_PING  = OpCode(2) // ping指令
    OPCODE_PONG  = OpCode(3) // pong指令
    OPCODE_DATA  = OpCode(4) // 数据消息
    OPCODE_TURN  = OpCode(5) // 切换线路 |1byte|string bytes| => |数据类型|连接地址| => |0|1.1.1.1:53|
    OPCODE_NOOP  = OpCode(6) // 链路踢出
    OPCODE_REQ   = OpCode(7) // 发送请求
    OPCODE_RESP  = OpCode(8) // 收到回复
)

type Packet struct {
    OpCode  OpCode
    Data    []byte
    DataLen uint32
}

func WriteMessage(conn net.Conn, op OpCode, data []byte) (int, error) {
    var (
        buff = make([]byte, 4)
    )
    if data == nil {
        data = []byte{}
    }

    binary.BigEndian.PutUint32(buff, uint32(len(data)))
    buff[0] = uint8(op) // data flag
    return conn.Write(bytes.Join([][]byte{buff, data}, []byte{}))
}

func ReadMessage(conn net.Conn) (*Packet, error) {
    var (
        buff = make([]byte, 4)
        err  error
        pkg  Packet
    )
    _, err = readFullWithTimeout(conn, buff, time.Second*10)
    if err != nil {
        return &pkg, err
    }
    op := OpCode(buff[0])

    buff[0] = 0
    bodyLen := binary.BigEndian.Uint32(buff)
    // limit size
    if bodyLen > 16*1024*1024 {
        // todo Too large Message
        return &pkg, ErrorMessageTooLarge
    }
    buff = make([]byte, bodyLen)
    n, err := readFullWithTimeout(conn, buff, time.Second*10)
    pkg.Data = buff
    pkg.DataLen = uint32(n)
    pkg.OpCode = op
    return &pkg, err
}

func readFullWithTimeout(conn net.Conn, buffer []byte, timeout time.Duration) (n int, err error) {
    if timeout == 0 {
        return io.ReadFull(conn, buffer)
    } else {
        if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
            return 0, err
        }
        return readAtLeast(conn, buffer, len(buffer))
    }
}

func readAtLeast(conn net.Conn, buf []byte, min int) (n int, err error) {
    if len(buf) < min {
        return 0, io.ErrShortBuffer
    }
    for n < min && err == nil {
        var nn int
        nn, err = conn.Read(buf[n:])
        n += nn
    }
    if n >= min {
        err = nil
    } else if n > 0 && err == io.EOF {
        err = io.ErrUnexpectedEOF
    }
    return
}
