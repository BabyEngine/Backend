
local rooms = {}
local roomId = 100
local playerRoomMap = {}

function NewRoom()
    local self = {
        id = roomId,
        players = {}
    }
    roomId = roomId + 1
    rooms[self.id] = self

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
        local lastRoom = playerRoomMap[conn]
        if lastRoom then
            if lastRoom == self then
                print('已经加入此房间了')
            else
                print('存在未完成的房间')
            end
            return
        end
        -- 新建玩家
        local player = GetPlayer(conn)
        player.room = self
        player.id = pid
        self.players[player.id] = player
        playerRoomMap[conn] = self
        local tmp = GetMyRoom(conn)
        print('映射玩家', conn, self, tmp)
    end

    function self.LeaveRoom( conn )
        local player = GetPlayer(conn)
        print('离开??', player)
        if player then
            self.players[player.id] = nil
            playerRoomMap[conn] = nil
            return player
        end
    end

    return self
end

function ListRoom()
    return rooms
end

function JoinRoom( roomId, playerId, conn )
    local room = rooms[roomId]
    if room then
        room.JoinRoom( roomId, playerId, conn )
        return room
    end
end

function LeaveRoom( conn )
    local room = GetMyRoom(conn)
    if room then
        return room.LeaveRoom( conn )
    end
end

function GetMyRoom( conn )
    return playerRoomMap[conn]
end

function onUpdateCheck( ... )
    for k,v in pairs(rooms) do
        v.Update()
    end
end

Looper.AddUpdate(onUpdateCheck)