function NewPlayerActor(client)
    local self = {
        conn = client,
        messageHandler = {},
        isRelease = false,
    }

    function self.Release()
        print('玩家离开了', self.token)
        self.isRelease = true
    end
    function self.OnMessage( data )
        local name, obj = GameMsg.decode(data)
        local cb = self.messageHandler[name]
        if cb then
            cb(obj)
        end
    end
    function self.OnRequest( data )

    end

    self.messageHandler['game.login'] = function(msg)
        print('欢迎登录', table.tostring(msg))
        self.isAuth = true
        self.token = msg.token
        netd.OnAuth(self.conn, self)
        local pushMsg
        pushMsg = function ()
            if self.isRelease then return end
            local msg = GameMsg.encode('game.login', {token='狗子'})
            self.conn.Send(msg)
            Looper.AfterFunc(2, pushMsg)
            print('going to send')
        end
        pushMsg()
    end

    return self
end