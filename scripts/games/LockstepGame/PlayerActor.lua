local players = {}

function NewPlayerActor(conn)
    local self = {
        id = id,
        conn = conn,
    }

    function self.OnPlayingMessage(msg)
        print('player playing message:', table.tostring(msg))
    end
    players[conn] = self
    return self
end

function GetPlayer( conn )
    return players[conn]
end