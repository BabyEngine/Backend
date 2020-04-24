net = net or {}

function net.NewHTTPServer(address, options)
    local self = {}
    options = options or {}
    local ptr
    local function onNew(conn)
        local r = json.decode(BabyEngine.Net.RunCmd(ptr, conn, 'get_all_query', '')) or {}
        if r and r.Body then
            r.Body = base64.decode(r.Body)
        end
        local client = {}
        function client.Close()
            BabyEngine.Net.Close(ptr, conn)
        end
        function client.Send( httpStatusCode, data )
            if not data then
                data = httpStatusCode
                httpStatusCode = 200
            end
            BabyEngine.Net.SendRaw(ptr, conn, tostring(httpStatusCode), 0, data)
        end
        if self.Serve then
            self.Serve(client, r)
        end
    end

    function self.Start( )
        ptr = BabyEngine.Net.Start('http', address, {tag=tag, ssl_key=options.ssl_key, ssl_cert=options.ssl_cert})
        BabyEngine.Net.Bind(ptr, "new",  onNew)
    end

    function self.Stop()
        if ptr then
            BabyEngine.Net.Stop(ptr)
            ptr = nil
        end
    end
    return self
end