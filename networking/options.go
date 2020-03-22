package networking

import (
    "context"
    "time"
)

type OptionFunc func(options *Options)
type Options struct {
    Type    string
    Tag     string
    Address string
    Handler ClientHandler
    TTL     time.Duration
    Ctx     context.Context
}

func WithType(t string) OptionFunc {
    return func(options *Options) {
        options.Type = t
    }
}

func WithTag(tag string) OptionFunc {
    return func(options *Options) {
        options.Tag = tag
    }
}

func WithAddress(address string) OptionFunc {
    return func(options *Options) {
        options.Address = address
    }
}

func WithHandler(handler ClientHandler) OptionFunc {
    return func(options *Options) {
        options.Handler = handler
    }
}

func WithContext(c context.Context) OptionFunc {
    return func(options *Options) {
        options.Ctx = c
    }
}

func DefaultOptions() *Options {
    opts := &Options{}
    opts.TTL = time.Second * 30
    opts.Ctx = context.TODO()
    return opts
}
