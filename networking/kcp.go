package networking


func mKCPListenAndServe(opts *Options) error {
    var (
        server mKCPServer
    )
    server.opts = opts
    server.clients = make(map[int64]*mKCPClient)
    return server.Serve(opts.Address)
}
