package core

import (
    "bufio"
    "bytes"
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "crypto/sha512"
    "encoding/hex"
    "github.com/DGHeroin/golua/lua"
    "hash"
    "io"
    "os"
)

func initModHash(L *lua.State) {
    L.GetGlobal("BabyEngine")
    L.PushString("Hash")
    {
        // 创建子表
        L.CreateTable(0, 1)

        L.PushString("HashFile")
        L.PushGoFunction(gHashCalcFile)
        L.SetTable(-3)

        L.PushString("HashString")
        L.PushGoFunction(gHashCalcString)
        L.SetTable(-3)
    }
    L.SetTable(-3)
}
func gHashCalcString(L *lua.State) int {
    var (
        reader io.Reader
        hashTypes = [] string{}
        writer = []io.Writer{}
        hashes = []hash.Hash{}
    )

    n := L.GetTop()
    str := L.ToString(1)
    reader = bytes.NewBufferString(str)
    for i := 2; i <= n; i++ {
        h := L.CheckString(i)
        hashTypes = append(hashTypes, h)
    }

    for _, v := range hashTypes {
        switch v {
        case "md5":
            hs := md5.New()
            writer = append(writer, hs)
            hashes = append(hashes, hs)
        case "sha1":
            hs := sha1.New()
            writer = append(writer, hs)
            hashes = append(hashes, hs)
        case "sha256":
            hs := sha256.New()
            writer = append(writer, hs)
            hashes = append(hashes, hs)
        case "sha512":
            hs := sha512.New()
            writer = append(writer, hs)
            hashes = append(hashes, hs)
        }
    }

    if err := CalculateHashes(reader, writer...); err == nil {
        L.PushBoolean(true)
        for _, h := range hashes {
            result := hex.EncodeToString(h.Sum(nil))
            L.PushString(result)
        }
        return 1 + len(hashes)
    }

    L.PushBoolean(false)
    return 1
}
func gHashCalcFile(L *lua.State) int {
    var (
        reader io.ReadCloser
        hashTypes = [] string{}
        writer = []io.Writer{}
        hashes = []hash.Hash{}
        err error
    )

    n := L.GetTop()
    if reader, err = os.Open(L.ToString(1)); err != nil {
        L.PushBoolean(false)
        return 1
    }
    defer reader.Close()
    for i := 2; i <= n; i++ {
        hashTypes = append(hashTypes, L.ToString(i))
    }

    for _, v := range hashTypes {
        switch v {
        case "md5":
            hs := md5.New()
            writer = append(writer, hs)
            hashes = append(hashes, hs)
        case "sha1":
            hs := sha1.New()
            writer = append(writer, hs)
            hashes = append(hashes, hs)
        case "sha256":
            hs := sha256.New()
            writer = append(writer, hs)
            hashes = append(hashes, hs)
        case "sha512":
            hs := sha512.New()
            writer = append(writer, hs)
            hashes = append(hashes, hs)
        }
    }

    if err := CalculateHashes(reader, writer...); err == nil {
        L.PushBoolean(true)
        for _, h := range hashes {
            result := hex.EncodeToString(h.Sum(nil))
            L.PushString(result)
        }
        return 1 + len(hashes)
    }

    L.PushBoolean(false)
    return 1
}

type HashInfo struct {
    Md5    string `json:"md5"`
    Sha1   string `json:"sha1"`
    Sha256 string `json:"sha256"`
    Sha512 string `json:"sha512"`
}

type HashCalculator interface {
    io.Writer
    Sum(b []byte) []byte
    Reset()
    Size() int
    BlockSize() int
}

func CalculateHashes(rd io.Reader, writer ... io.Writer) error {
    pageSize := os.Getpagesize()
    reader := bufio.NewReaderSize(rd, pageSize)
    multiWriter := io.MultiWriter(writer...)
    _, err := io.Copy(multiWriter, reader)
    return err
}