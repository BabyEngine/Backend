package Debug

import "log"

var (
    Log = log.Println
    Logf = log.Printf
    LogIf = logIf
    LogIff = logIff
)

func init()  {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func logIf(cond bool, v ...interface{}) {
    if cond {
        Log(v...)
    }
}

func logIff(cond bool, format string, v ...interface{}) {
    if cond {
        Logf(format, v...)
    }
}