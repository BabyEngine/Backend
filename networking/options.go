package networking

import (
    "context"
    "time"
)

type OptionFunc func(options *Options)
type Options struct {
    Type      string
    Tag       string
    Address   string
    Handler   ClientHandler
    TTL       time.Duration
    Ctx       context.Context
    IsRawMode bool
    TLSKey    string
    TLSCert   string
    TLSEnable bool
    Args map[string]string
}

func WithArgs(args map[string]string) OptionFunc {
    return func(options *Options) {
        options.Args = args
    }
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

func WithRawMode(b bool) OptionFunc {
    return func(options *Options) {
        options.IsRawMode = b
    }
}

func WithTLS(key, cert string) OptionFunc {
    return func(options *Options) {
        options.TLSKey = key
        options.TLSCert = cert
        if options.TLSKey != "" && options.TLSCert != "" {
            options.TLSEnable = true
        } else {
            options.TLSEnable = false
        }
    }
}

func DefaultOptions() *Options {
    opts := &Options{}
    opts.TTL = time.Second * 30
    opts.Ctx = context.TODO()
    return opts
}
