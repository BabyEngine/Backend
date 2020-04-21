net = net or {}

function net.NewHTTPServer(address)
    local self = {}
    local ptr
    local function onNew(conn)
        -- print('http cli new', conn)

        -- print(BabyEngine.Net.RunCmd(ptr, conn, 'header', '*'))
        -- print('body:', BabyEngine.Net.RunCmd(ptr, conn, 'body', '5'))
        -- print(BabyEngine.Net.RunCmd(ptr, conn, 'host', ''))
        -- print(BabyEngine.Net.RunCmd(ptr, conn, 'method', ''))
        -- print(BabyEngine.Net.RunCmd(ptr, conn, 'query', ''))

        -- BabyEngine.Net.RunCmd(ptr, conn, 'set_header', 'custom_key', 'custom_value')
        -- BabyEngine.Net.Send(ptr, conn, 'hello!!!!')
        -- BabyEngine.Net.Close(ptr, conn)
        -- print(BabyEngine.Net.RunCmd(ptr, conn, 'get_all_query', ''))
        local r = json.decode(BabyEngine.Net.RunCmd(ptr, conn, 'get_all_query', '')) or {}
        if r and r.Body then
            r.Body = base64.decode(r.Body)
        end
        local client = {}
        function client.Close()
            BabyEngine.Net.Close(ptr, conn)
        end
        function client.Send( data )
            BabyEngine.Net.Send(ptr, conn, data)
        end
        if self.Serve then
            self.Serve(client, r)
        end

    end
    local function onClose(conn)
        --print('http cli close', conn)
    end
    local function onError(conn, err)
        --print('http cli on error', conn)
    end

    function self.Start( )
        ptr = BabyEngine.Net.Start('http', address, {tag=tag})
        BabyEngine.Net.Bind(ptr, "new",  onNew)
        BabyEngine.Net.Bind(ptr, "close",  onClose)
        BabyEngine.Net.Bind(ptr, "error",  onError)
    end

    function self.Stop()
        if ptr then
            BabyEngine.Net.Stop(ptr)
            ptr = nil
        end
    end
    return self
end