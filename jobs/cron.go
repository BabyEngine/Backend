package jobs

import (
    "github.com/robfig/cron/v3"
)

type Context struct {
    cron *cron.Cron
}

func NewCron() *Context {
    c := &Context{}
    c.cron = cron.New(cron.WithSeconds())
    c.cron.Start()
    return c
}

func (c *Context) Add(spec string, cb func()) (int, error) {

    id, err := c.cron.AddFunc(spec, cb)
    if err != nil {
        return -1, err
    }
    return int(id), nil
}

func (c *Context) Remove(id int) {
    c.cron.Remove(cron.EntryID(id))
}

func (c *Context) Stop()  {
    c.cron.Stop()
}