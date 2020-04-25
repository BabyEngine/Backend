rpc = rpc or {}

function rpc.NewServer()
    local self = {}

    local function handleFunc(action, data, respFunc)
        if self.OnRequest then
            self.OnRequest(action, data, function ( code, data )
                if type(code) ~= 'number' then
                    respFunc(-1, '')
                end
                if type(data) ~= 'string' then
                    respFunc(-2, '')
                end
                respFunc(code, data)
            end)
        end
    end

    function self.Listen( address )
        BabyEngine.RPC.NewServer(address, handleFunc)
    end

    return self
end

function rpc.NewClient()
    local self = {}
    local clientPtr = nil
    function self.Connect( address )
        local client, err = BabyEngine.RPC.NewClient('127.0.0.1:9981')
        if err ~= nil then
            return err
        end
        clientPtr = client
    end
    function self.Call( action, data, cb )
        BabyEngine.RPC.Call(clientPtr, action, data, function(code, data, err)
            if type(cb) == 'function' then
                cb(code, data, err)
            end
        end)
    end
    return self
end