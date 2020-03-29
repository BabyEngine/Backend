
local rooms = {}
local roomId = 100
function NewRoom()
    local self = {
        id = roomId,
        players = {}
    }
    roomId = roomId + 1
    table.insert(rooms, self)

    function self.Update()

    end

    function self.GetPlayers( )
        local result = {}
        for k,v in pairs(self.players) do
            table.insert(result, v.id)
        end
        return result
    end

    function self.JoinRoom( rid, pid, conn )
        for i,v in ipairs(self.players) do
            if v.id == pid then
                print('已经在房间了')
                return
            end
        end
        -- 新建玩家
        local player = NewPlayerActor(pid, conn)
        table.insert(self.players, player)
    end
    return self
end

function ListRoom()
    return rooms
end

function JoinRoom( roomId, playerId, conn )
    for k,v in pairs(rooms) do
        if v.id == roomId then
            v.JoinRoom( roomId, playerId, conn )
            return v
        end
    end
end

function onUpdateCheck( ... )
    for k,v in pairs(rooms) do
        v.Update()
    end
end

Looper.AddUpdate(onUpdateCheck)