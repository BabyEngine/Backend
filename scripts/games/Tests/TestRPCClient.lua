local client = jsonrpc.NewClient()
local err = client.Connect('http://127.0.0.1:9981')
if err ~= nil then
    print(err)
    return
end

Looper.AddTick(1, function ()
    client.Call('payment.list', 'user-id', function ( code, data, err )
        print(code, data, err)
    end)
end)

