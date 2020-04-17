package debugging

import "log"

var (
    Log = log.Println
    Logf = log.Printf
    LogIf = logIf
    LogIff = logIff

    customPrefix = ""
)

func init()  {
    log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
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

func Prefix(prefix string) {
    customPrefix = prefix
    log.SetPrefix(customPrefix)
}