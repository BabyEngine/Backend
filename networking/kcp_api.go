package networking

func ListenAndServe(options ...OptionFunc) error {
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
    case "ws":
        var server mWebsocketServer
        server.opts = opts
        server.Init()
        return server.Serve(opts.Address)
    case "http":
        var server mHTTPServer
        server.opts = opts
        server.Init()
        return server.Serve(opts.Address)
    case "socket.io":
        var server mSocketIOServer
        server.opts = opts
        server.Init()
        return server.Serve(opts.Address)
    }
    return ErrorOptionsInvalid
}
