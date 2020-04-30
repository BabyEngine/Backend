local server = rpc.NewServer('jsonrpc')
server.OnRequest = function(action, data, respFunc)
    respFunc(0, 'jsonrpc reply msg')
end

server.Listen(':9981')


local server2 = rpc.NewServer('grpc')
server2.OnRequest = function(action, data, respFunc)
    respFunc(0, 'grpc reply msg')
end

server2.Listen(':9982')