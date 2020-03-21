
local test_code = [[
local db = BabyEngine.KV.Open('my.db')
print(db)

local err, ok = BabyEngine.KV.Put(db, 'users', 'mytoken', 'iqoo')
print('写数据', err, '=>', ok)

BabyEngine.KV.Get(db, 'users', 'mytoken', function(val, err)
    print('读数据', err, '=>', val)
end)
]]
