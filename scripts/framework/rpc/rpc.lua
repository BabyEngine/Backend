rpc = rpc or {}

function rpc.NewServer(t)
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
        BabyEngine.RPC.NewServer(t, address, handleFunc)
    end

    return self
end

function rpc.NewClient(t, address)
    local self = {}
    local clientPtr = nil
    clientPtr, err = BabyEngine.RPC.NewClient(t, address)
    function self.Connect()
        -- return ok, errStr
        return BabyEngine.RPC.Connect(clientPtr)
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