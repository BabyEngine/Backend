local cron = Jobs.NewCron()
local id, err = cron.Add('@every 1s', function ( ... )
    print('1 seconds job')
end)

cron.Add('@every 2s', function ( ... )
    print('2 seconds job')
end)


Looper.AfterFunc(3, function( ... )
    cron.Remove(id)
end)

Looper.AfterFunc(7, function( ... )
    cron.Stop()
    cron = nil
end)