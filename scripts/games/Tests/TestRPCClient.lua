local client = rpc.NewClient('jsonrpc', 'http://127.0.0.1:9981')
local ok, err = client.Connect()
if err ~= nil then
    print(err)
    return
end

Looper.AddTick(1, function ()
    client.Call('payment.list', 'user-id', function ( code, data, err )
        print(code, data, err)
    end)
end)

local client2 = rpc.NewClient('grpc', '127.0.0.1:9982')
local ok, err = client2.Connect()
if err ~= nil then
    print(err)
    return
end

Looper.AddTick(1, function ()
    client2.Call('payment.list.grpc', 'user-id', function ( code, data, err )
        print(code, data, err)
    end)
end)

