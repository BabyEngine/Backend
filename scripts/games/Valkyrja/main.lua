print('女武神启动')

require('games.Valkyrja.errors.errors')
require('games.Valkyrja.net.pbdef')
require('games.Valkyrja.net.GameMsg')
require('games.Valkyrja.net.service')
require('games.Valkyrja.misc.db')
require('games.Valkyrja.player.PlayerActor')

-- 设置update频率
BabyEngine.App.SetFPS(60)
-- 监听服务器端口
NetService.startGameServer(":8087")