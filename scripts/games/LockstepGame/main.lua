require('games.LockstepGame.pb')
require('games.LockstepGame.Server')
require('games.LockstepGame.PlayerActor')
require('games.LockstepGame.Room')

-- 设置update频率
BabyEngine.App.SetFPS(60)
-- 监听服务器端口
NetService.startGameServer(":8087")

