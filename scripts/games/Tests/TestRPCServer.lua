local server = rpc.NewServer()
server.OnRequest = function(action, data, respFunc)
    respFunc(0, 'reply msg')
end

server.Listen(':9981')