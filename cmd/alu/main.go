package main

import (
    "archive/tar"
    "fmt"
    "github.com/BabyEngine/Backend/logger"
    "github.com/BabyEngine/Backend/game"
    "github.com/BabyEngine/Backend/hotzone"
    "github.com/DGHeroin/golua/lua"
    "io"
    "io/ioutil"

    "net/http"
    "os"
    "path/filepath"
    "strings"
)

var (
    app *game.Application
    Version string = "v00.00.00"
)

func main() {
    if len(os.Args) == 1 {
        printHelp()
        return
    }
    cmd := os.Args[1]
    switch cmd {
    case "version":
        fmt.Println(Version)
    case "run":
        runLuaApp()
    case "get":
        runGetAction()
    case "help":
        printHelp()
    default:
        runLuaApp()
    }
}

func printHelp() {
    usage:=`
    version | print current version
    run main.lua | run a lua file
    get framework.tar | get package form github repo
`
    fmt.Println(usage)
}

func runLuaApp() {
    app = game.NewApp()
    if os.Getenv("HotReload") == "true" {
       go hotzone.EnableHotRestart(app, runLuaApp)
    }
    L := lua.NewState()
    L.OpenLibs()
    L.OpenGoLibs()
    defer L.Close()

    if len(os.Args) == 1 {
       logger.Warn("no input file")
       return
    }
    var file string
    if len(os.Args) == 2 {
        file = os.Args[1]
    } else if len(os.Args) == 3 {
        file = os.Args[2]
    }

    app.Init(L)

    if err := L.DoFile(file); err != nil {
        logger.Warn(err)
    }
    app.Start()
}

func runGetAction() {
    if len(os.Args) < 3 {
        logger.Warn("args error")
        return
    }
    urlTemplate := []string{
        "https://github.com/BabyEngine/Backend/releases/latest/download/%s",
        "https://github.com/BabyEngine/Backend/releases/download/%s/%s",
    }
    var (
        tokens []string
    )
    what := os.Args[2]

    if strings.HasPrefix(what, "github.com") {
        urlTemplate = []string{
            "https://"+what+"/releases/latest/download/%s",
            "https://"+what+"/releases/download/%s/%s",
        }
        what = os.Args[3]
        tokens = strings.Split(what, "@")
        if len(tokens) < 1 {
            logger.Warn("package not found")
        }
    } else {
        tokens = strings.Split(what, "@")
        if len(tokens) < 1 {
            logger.Warn("package not found")
        }
    }

    var (
        packageName string
        version string
        downloadUrl string
    )
    packageName = tokens[0]

    if len(tokens) == 2 {
        version = tokens[1]
    }

    //logger.Debugf("%s %s", packageName, version)
    if version == "" { // download latest version
        downloadUrl = fmt.Sprintf(urlTemplate[0], packageName)
    } else { // download version
        downloadUrl = fmt.Sprintf(urlTemplate[1], version, packageName)
    }
    resp, err := http.Get(downloadUrl)
    if err != nil {
       logger.Warnf("download error", err)
       return
    }
    if resp.StatusCode != http.StatusOK {
       logger.Warnf("%s %s", resp.Status, downloadUrl)
       return
    }

    defer resp.Body.Close()
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      logger.Warnf("read data error:%v", err)
      return
    }
    if data == nil {
      return
    }

    dname := ".tmp/bbe"
    defer os.RemoveAll(".tmp")
    if _, err := os.Stat(dname); os.IsNotExist(err) {
        os.MkdirAll(dname, os.ModePerm)
    }
    fname := filepath.Join(dname, packageName)
    err = ioutil.WriteFile(fname, data, os.ModePerm)
    if err != nil {
      logger.Warnf("write file error:%v", err)
      return
    }

    if err := untar(fname, "./"); err != nil {
        logger.Warnf("%s", err)
    }
}

func untar(archive, dst string) error {
    // 打开准备解压的 tar 包
    fr, err := os.Open(archive)
    if err != nil {
        return err
    }
    defer fr.Close()
    tr := tar.NewReader(fr)
    for {
        hdr, err := tr.Next()
        switch {
        case err == io.EOF:
            return nil
        case err != nil:
            return err
        case hdr == nil:
            continue
        }
        dstFileDir := filepath.Join(dst, hdr.Name)
        switch hdr.Typeflag {
        case tar.TypeDir:
            if b := ExistDir(dstFileDir); !b {
                if err := os.MkdirAll(dstFileDir, 0775); err != nil {
                    return err
                }
            }
        case tar.TypeReg:
            file, err := os.OpenFile(dstFileDir, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
            if err != nil {
                fmt.Println(err)
                return err
            }
            n, err := io.Copy(file, tr)
            if err != nil {
                fmt.Println(err)
                return err
            }
            fmt.Printf("成功解压： %s , 共处理了 %d 个字符\n", dstFileDir, n)
            _ = file.Close()
        }
    }
}
// 判断目录是否存在
func ExistDir(dirname string) bool {
    fi, err := os.Stat(dirname)
    return (err == nil || os.IsExist(err)) && fi.IsDir()
}
