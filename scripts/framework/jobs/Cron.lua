Jobs = Jobs or {}

function Jobs.NewCron()
    local self = {}
    local corn = BabyEngine.Cron.New()
    function self.Add( sepc, cb )
        return BabyEngine.Cron.Add(corn, sepc, cb)
    end

    function self.Remove( id )
        return BabyEngine.Cron.Remove(corn, id)
    end

    function self.Stop( id )
        return BabyEngine.Cron.Stop(corn)
    end
    return self
end