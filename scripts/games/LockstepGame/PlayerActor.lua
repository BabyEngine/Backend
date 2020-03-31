function NewPlayerActor(id, client)
    local self = {
        id = id,
        conn = client,
    }

    function self.OnPlayingMessage(msg)
        print('player playing message:', table.tostring(msg))
    end

    return self
end
