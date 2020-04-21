package networking

import "C"
import (
    "encoding/json"
    "errors"
    "github.com/BabyEngine/Backend/debugging"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"
)

type mHTTPClient struct {
    server     *mHTTPServer
    id         int64
    w          http.ResponseWriter
    r          *http.Request
    opts       *Options
    stopChan   chan interface{}
    isStopRead bool
    isStop     bool
}

func (c *mHTTPClient) init() {
    c.stopChan = make(chan interface{}, 1)
    c.isStopRead = false
}

func (c *mHTTPClient) SendData(data []byte) error {
    _, err := c.w.Write(data)
    return err
}

func (c *mHTTPClient) SendRaw(op OpCode, data []byte) error {
    return errors.New("UnImpl")
}

func (c *mHTTPClient) Close() {
    if c.isStopRead {
        return
    }
    if c.isStop {
        return
    }
    c.isStop = true

    go func() {
        defer c.r.Body.Close()
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

func (c *mHTTPClient) Id() int64 {
    return c.id
}
func (c *mHTTPClient) SetId(id int64) {
    c.id = id
}

func (c *mHTTPClient) Serve() {
    c.opts.Handler.OnNew(c)

EXITLOOP:
    for {
        select {
        case <-c.stopChan:
            break EXITLOOP
        }
    }
}

func (c *mHTTPClient) RunCmd(cmd string, args []string) string {
    if len(args) < 1 {
        return ""
    }
    key := args[0]
    switch cmd {
    case "header":
        if key == "*" {
            if js, err := json.Marshal(c.r.Header); err == nil {
                return string(js)
            } else {
                debugging.Log(err)
            }
        }
        return c.r.Header.Get(key)
    case "body":
        if key == "*" {
            if data, err := ioutil.ReadAll(c.r.Body); err == nil {
                return string(data)
            } else {
                debugging.Log(err)
            }
        }
        if n, err := strconv.Atoi(key); err == nil {
            buff := make([]byte, n)
            if n, err := c.r.Body.Read(buff); err == nil {
                return string(buff[:n])
            } else {
                debugging.Log(err)
            }
        } else {
            debugging.Log(err)
        }
    case "method":
        return c.r.Method
    case "host":
        return c.r.Host
    case "query":
        return c.r.URL.RawQuery
    case "set_header":
        if len(args) < 2 {
            return ""
        }
        c.w.Header().Set(args[1], args[2])
    case "get_all_query":
        type QueryInfo struct {
            Header http.Header
            Query  string
            Method string
            Host   string
            Body   []byte
            Path   string
        }
        body, err := ioutil.ReadAll(c.r.Body)
        if err != nil {
        }
        data, _ := json.Marshal(&QueryInfo{
            Header: c.r.Header,
            Query:  c.r.URL.RawQuery,
            Method: c.r.Method,
            Host:   c.r.Host,
            Body:   body,
            Path:   c.r.URL.Path,
        })
        return string(data)
    }
    return ""
}
