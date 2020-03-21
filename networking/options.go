package networking

import "context"

type OptionFunc func(options *Options)
type Options struct {
    Type    string
    Tag     string
    Address string
    Handler ClientHandler
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

func (o *Options) Valid() error {
    if o.Tag == "" || o.Handler == nil || o.Address == "" || o.Type == "" {
        return ErrorOptionsInvalid
    }
    return nil
}
