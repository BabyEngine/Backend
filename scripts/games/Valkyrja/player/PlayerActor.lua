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
        local name, obj = GameMsg.decode(data)
        local cb = self.messageHandler[name]
        if cb then
            cb(obj, respFunc)
            return
        end
        print('请求什么', name)
        respFunc(GameMsg.encode('rsp.common', {code=-1, msg=errors.ACTION_NOT_FOUND}))
        return nil
    end

    self.messageHandler['req.login'] = function(msg, respFunc)
        self.isAuth = true
        self.token = msg.token
        GetOrCreatePlayerInfo(self.token, function(data, err)
            if err ~= nil then -- 认证失败
                print('认证失败', err)
                respFunc(GameMsg.encode('rsp.common', {code=-2, msg=errors.AUTH_FAILED}))
                return
            end
            -- 认证通过
            NetService.OnAuth(self.conn, self)
            self.playerInfo = data
            respFunc(GameMsg.encode('rsp.common', {code=0, msg=errors.OK}))
        end)
    end

    self.messageHandler['req.battle.start'] = function(msg, respFunc)
        print('开始战斗', self.playerInfo)
        respFunc(GameMsg.encode('rsp.common', {code=0, msg=errors.OK}))
        print('xxxx')
    end

    return self
end

function GetOrCreatePlayerInfo( token, cb )
    local tokens = string.split(token, ':')
    if not tokens then
         cb(nil, errors.TOKEN_INVALID)
    end
    local srouce = tokens[1] -- 来源 guest/regular
    local userId = tokens[2]
    local auth   = tokens[3]
    print('login token:', srouce, userId, auth)
    DB.GetUser(userId, function ( data, err )
        if err ~= nil then
            print('发生错误', err)
            cb(nil, err)
            return
        end
        if data == nil then
            -- 玩家数据不存在, 创建一条新的数据
            result = {userId=userId, auth =auth, source=source}
            DB.SaveUser(userId, json.encode(result))
        else
            result = json.decode(data)
            if result.auth ~= auth then --认证失败
                cb(result, errors.AUTH_FAILED)
                return
            end
        end
        cb(result, err)
    end)
end