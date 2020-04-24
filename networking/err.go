package networking

import "errors"

var (
    ErrorTimeout         = errors.New("timeout")
    ErrorMessageTooLarge = errors.New("message body too large")
    ErrorOptionsInvalid  = errors.New("opts invalid")
    ErrorAPINotSupport   = errors.New("api not support")
)
