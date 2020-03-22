DB = {}
local db = nil
function DB.GetUser( userId, cb )
    BabyEngine.KV.Get(db, 'users', userId, cb)
end

function DB.SaveUser( userId, data )
    BabyEngine.KV.Put(db, 'users', userId, data)
end

db = BabyEngine.KV.Open('my.db')