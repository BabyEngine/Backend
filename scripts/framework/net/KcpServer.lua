net = net or {}

function net.NewKCPBinaryServer(address, tag)
    local self = {
    }
    local clients = {}
    local ptr = nil
    local function onNew(conn)
        local cli = {
            conn = conn,
        }
        function cli.Send(data)
            SendNetData(ptr, cli.conn, data)
        end
        function cli.Close()
            CloseNetClient(ptr, cli)
        end
        function cli.Redirect(address)
            RedirectNetClient(ptr, cli, address)
        end
        clients[conn] = cli
        if not self.OnNew then return end
        self.OnNew(cli)
    end
    local function onClose(conn)
        local cli = clients[conn]
        if not cli then return end
        clients[conn] = nil
        if not self.OnClose then return end
        self.OnClose(cli)
    end
    local function onError(conn, err)
        local cli = clients[conn]
        if not cli then return end
        if not self.OnError then return end
        self.OnError(cli, err)
    end
    local function onData(conn, data)
        local cli = clients[conn]
        if not cli then return end
        if not self.OnData then return end
        self.OnData(cli, data)
    end
    local n = 0
    local function onRequest(conn, data)
        local cli = clients[conn]
        if not cli then return end
        if not self.OnRequest then return end
        return self.OnRequest(cli, data)
    end
    function self.Start( )
        ptr = StartNetServer('kcp', address, tag)
        BindNetServer(ptr, "new",  onNew)
        BindNetServer(ptr, "close",  onClose)
        BindNetServer(ptr, "error",  onError)
        BindNetServer(ptr, "data",  onData)
        BindNetServer(ptr, "request",  onRequest)
    end

    function self.Stop()
        if ptr then
            StopNetServer(ptr)
            ptr = nil
        end
    end
    return self
end
