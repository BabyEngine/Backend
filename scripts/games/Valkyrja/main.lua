print('女武神启动')

BabyEngine.App.SetFPS(1)

require('games.Valkyrja.netd.GameMsg')
require('games.Valkyrja.player.PlayerActor')
require('games.Valkyrja.netd.service')

netd.startGameServer(":8087")

print(table.tostring(BabyEngine))

