net = net or {}

function net.NewSocketIOServer(address, options)
    local self = {}
    options = options or {}
    local clients = {}
    local ptr = nil
    local function onNew(conn)
        local cli = {
            conn = conn,
        }

        function cli.Emit(event, data)
            BabyEngine.Net.SendRaw(ptr, cli.conn, event, 0, data)
        end
        function cli.Close()
            BabyEngine.Net.Close(ptr, cli.conn)
        end
        function cli.Redirect(address)
            BabyEngine.Net.Redirect(ptr, cli.conn, address)
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
    local function onRequest(conn, data, respFunc)
        local cli = clients[conn]
        if not cli then return end
        if not self.OnRequest then return end
        self.OnRequest(cli, data, respFunc)
    end

    function self.Start( )
        ptr = BabyEngine.Net.Start('socket.io', address,
            {raw='true',
                eventName = options.eventName,
            ssl_key=options.ssl_key,
            ssl_cert=options.ssl_cert})
        BabyEngine.Net.Bind(ptr, "new",  onNew)
        BabyEngine.Net.Bind(ptr, "close",  onClose)
        BabyEngine.Net.Bind(ptr, "error",  onError)
        BabyEngine.Net.Bind(ptr, "request",  onRequest)
    end

    function self.Stop()
        if ptr then
            BabyEngine.Net.Stop(ptr)
            ptr = nil
        end
    end
    return self
end