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
    function self.OnRequest( data, respFunc )
        print('req...', data)
        local name, obj = GameMsg.decode(data)
        local cb = self.messageHandler[name]
        if cb then
            cb(obj)
        end
        Looper.AfterFunc(2, function()
            print('稍后回复')
            respFunc(data)
        end)
        return nil
        --return data
    end

    self.messageHandler['game.login'] = function(msg)
        print('欢迎登录', table.tostring(msg))
        self.isAuth = true
        self.token = msg.token
        GetOrCreatePlayerInfo(self.token, function(data, err)
            if err ~= nil then -- 认证失败
                print('认证失败', err)
                return
            end
            -- 认证通过
            print('认证通过')
            NetService.OnAuth(self.conn, self)
            self.playerInfo = data
        end)
    end

    return self
end

function GetOrCreatePlayerInfo( token, cb )
    local tokens = string.split(token, ':')
    if not tokens then
         cb(nil, errors.TOKEN_INVALID)
    end
    local userId = tokens[1]
    local auth   = tokens[2]

    DB.GetUser(userId, function ( data, err )
        if err ~= nil then
            print('发生错误', err)
            cb(nil, err)
            return
        end
        if data == nil then
            -- 玩家数据不存在, 创建一条新的数据
            result = {userId=userId, auth =auth}
            DB.SaveUser(userId, json.encode(result))
        else
            result = json.decode(data)
            if result.auth ~= auth then --认证失败
                print('token 不一样', auth, '=>', result.auth)
                cb(result, errors.AUTH_FAILED)
                return
            end
        end
        print('获取玩家数据', table.tostring(result))
        cb(result, err)
    end)
end