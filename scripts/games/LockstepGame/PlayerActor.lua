function NewPlayerActor(id, client)
    local self = {
        id = id,
        conn = client,
    }


    return self
end
