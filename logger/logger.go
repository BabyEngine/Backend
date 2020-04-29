package logger

import (
    "fmt"
    "os"
    "runtime"
    "runtime/debug"
    "time"
)

var (
    Debug    func(v ...interface{})
    Debugf   func(format string, v ...interface{})
    DebugIf  func(cond bool, v ...interface{})
    DebugIff func(cond bool, format string, v ...interface{})
    Warn     func(v ...interface{})
    Warnf    func(format string, v ...interface{})
    WarnIf   func(cond bool, v ...interface{})
    WarnIff  func(cond bool, format string, v ...interface{})
    Error    func(v ...interface{})
    Errorf   func(format string, v ...interface{})
    ErrorIf  func(cond bool, v ...interface{})
    ErrorIff func(cond bool, format string, v ...interface{})
    Println  func(v ...interface{})
)

func init() {
    Debug = d_log
    Debugf = d_logf
    DebugIf = d_logIf
    DebugIff = d_logIff
    Warn = w_log
    Warnf = w_logf
    WarnIf = w_logIf
    WarnIff = w_logIff
    Error = e_log
    Errorf = e_logf
    ErrorIf = e_logIf
    ErrorIff = e_logIff
    Println = Debug
}
func shortFile() string {
    _, file, line, ok := runtime.Caller(3)
    if !ok {
        file = "???"
        line = 0
    } else {
        short := file
        for i := len(file) - 1; i > 0; i-- {
            if file[i] == '/' {
                short = file[i+1:]
                break
            }
        }
        file = short
    }

    return fmt.Sprintf("%v:%d", file, line)
}

func timestamp() string {
    return time.Now().UTC().Format(time.RFC3339Nano)
}

func writeLog(level int, v ...interface{}) {
    writeLogf(level, "%v", fmt.Sprint(v...))
}

func writeLogf(level int, format string, v ...interface{}) {
    if level == 1 {
        fmt.Fprintf(os.Stdout, "%s [debug] %s %v\n", timestamp(), shortFile(), fmt.Sprintf(format, v...))
    } else if level == 2 {
        fmt.Fprintf(os.Stdout, "%s [warn] %s %s", timestamp(), shortFile(), fmt.Sprintf(format, v...))
    } else if level == 3 {
        fmt.Fprintf(os.Stderr, "%s [error] %s %s\n%s", timestamp(), shortFile(), fmt.Sprintf(format, v...), string(debug.Stack()))
    }
}

func d_log(v ...interface{}) {
    writeLog(1, v...)
}
func d_logf(format string, v ...interface{}) {
    writeLogf(1, format, v...)
}
func d_logIf(cond bool, v ...interface{}) {
    if cond {
        writeLog(1, v...)
    }
}

func d_logIff(cond bool, format string, v ...interface{}) {
    if cond {
        writeLogf(1, format, v...)
    }
}
func w_log(v ...interface{}) {
    writeLog(2, v...)
}
func w_logf(format string, v ...interface{}) {
    writeLogf(2, format, v...)
}
func w_logIf(cond bool, v ...interface{}) {
    if cond {
        writeLog(2, v...)
    }
}

func w_logIff(cond bool, format string, v ...interface{}) {
    if cond {
        writeLogf(2, format, v...)
    }
}
func e_log(v ...interface{}) {
    writeLog(3, v...)
}
func e_logf(format string, v ...interface{}) {
    writeLogf(3, format, v...)
}
func e_logIf(cond bool, v ...interface{}) {
    if cond {
        writeLog(3, v...)
    }
}

func e_logIff(cond bool, format string, v ...interface{}) {
    if cond {
        writeLogf(3, format, v...)
    }
}
