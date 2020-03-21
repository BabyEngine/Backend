require('games.Valkyrja.netd.pbdef')

local pb = require "pb"
local protoc = require "protobuf/protoc"

if not pb.type"GameMessage" then
    protoc:load(PBDefine)
end

GameMsg = {}
local ParseMap = {}
function GameMsg.decode( data )
    local rs = pb.decode('GameMessage', data)
    if rs then
        if ParseMap[rs.action] then
            return rs.action, pb.decode(ParseMap[rs.action], rs.data)
        end
        return rs, nil
    end
    return
end


function GameMsg.encode( name, t )
    local obj = {
        action=name,
    }
    return pb.encode('GameMessage', obj)
end

ParseMap['game.login'] = 'RequestLoginMessage'