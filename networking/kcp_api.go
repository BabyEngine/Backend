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
    }
    return ErrorOptionsInvalid
}
