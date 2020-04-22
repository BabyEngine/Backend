--print('Init BabyEngine Backend')
require('framework.types')
require('framework.base64')
md5 = require('framework.md5')

require('framework.funcs')
-- require('framework.utils')
require('framework.looper')
require('framework.kv.db')

require('framework.net.KcpServer')
require('framework.net.WebsocketServer')
require('framework.net.HTTPServer')
require('framework.net.SocketIOServer')