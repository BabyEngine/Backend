package networking

type ClientHandler interface {
    OnNew(client Client)
    OnData(client Client, data []byte)
    OnClose(client Client)
    OnError(client Client, err error)
    OnRequest(client Client, data[]byte) []byte
    Stop()
    GetAllClient() []Client
    Mapping(key interface{}, client Client)
    GetClientByKey(key interface{}) Client
}

type Client interface {
    SendData(data []byte) error
    SendRawEvent(e string, op OpCode, data []byte) error
    Close()
    Id() int64
    SetId(id int64)
    RunCmd(action string, args[]string) string
    String() string
}