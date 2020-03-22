package networking

import "github.com/DGHeroin/golua/lua"

func newKCP(L *lua.State, address string, tag string) ClientHandler {
    // 主服务器
    h := &KCPGameServerHandler{}
    h.L = L
    h.Init()
    // 配置 Server
    go func() {
        if err := listenAndServe(
            WithType("kcp"),
            WithTag(tag),
            WithAddress(address),
            WithContext(h.ctx),
            WithHandler(h)); err != nil {
        }
    }()

    return h
}

func listenAndServe(options ...OptionFunc) error {
    opts := DefaultOptions()
    for _, cb := range options {
        cb(opts)
    }
    switch opts.Type {
    case "kcp":
        var server mKCPServer
        server.opts = opts
        server.Init()
        return server.Serve(opts.Address)
    }
    return ErrorOptionsInvalid
}
