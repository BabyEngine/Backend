print('欢迎来到游戏世界')

BabyEngine.App.SetFPS(60)

-- 启动NameServer
-- 客户端链接来后, 立刻把它分配到对应服务器上, 根据负载进行分流
function startNameServer( addr )
    local server = net.NewKCPBinaryServer(":8087", "NameServer")
    server.OnNew = function ( cli )
        print('新连接', table.tostring(cli))
        cli.Redirect('127.0.0.1:8088')
    end
    server.Start()
end

-- 游戏消息服务器
-- 处理客户端发到服务端的消息
function startGameServer( addr )
    local server = net.NewKCPBinaryServer(":8088", "GameServer")
    server.OnNew = function ( cli )
        print('新连接', table.tostring(cli))
    end

    server.OnClose = function ( cli )
        print('关闭连接', table.tostring(cli))
    end

    server.OnData = function ( cli, data )
        print('收到数据', table.tostring(cli), tostring(data))
    end

    server.OnRequest = function ( cli, data )
        return "Okay, Got it"
    end

    server.Start()
end


startNameServer()
startGameServer()

-- websocket服务器
function startWebsocketServer()
    local server = net.NewWebsocketBinaryServer(":8089", "websocket服务器")
    server.OnNew = function ( cli )
        print('新连接', cli)
    end

    server.OnClose = function ( cli )
        print('关闭连接', cli)
    end

    server.OnData = function ( cli, data )
        print('收到数据', cli, tostring(data))
        cli.Send('echo:'..tostring(data))
    end

    server.Start()
end

startWebsocketServer()

-- http 服务器
function startHTTPServer()
    local server = net.NewHTTPServer("", "http 服务器")
    server.Serve = function (cli, req)
        print('req', cli, table.tostring(req))
        cli.Send('服务 ok: ' .. md5.sumhexa('hello'))
        cli.Close()
    end

    server.Start()
end

startHTTPServer()