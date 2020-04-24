package game

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "errors"
    "github.com/DGHeroin/golua/lua"
    "io"
)

func initModCipher(L *lua.State) {
    L.GetGlobal("BabyEngine")
    L.PushString("Cipher")
    {
        // 创建子表
        L.CreateTable(0, 1)

        L.PushString("RSACreate")
        L.PushGoFunction(gRSACreate)
        L.SetTable(-3)

        L.PushString("RSAEncrypt")
        L.PushGoFunction(gRSAEncrypt)
        L.SetTable(-3)

        L.PushString("RSADecrypt")
        L.PushGoFunction(gRSADecrypt)
        L.SetTable(-3)

        L.PushString("AESEncrypt")
        L.PushGoFunction(gAESEncrypt)
        L.SetTable(-3)

        L.PushString("AESDecrypt")
        L.PushGoFunction(gAESDecrypt)
        L.SetTable(-3)
    }
    L.SetTable(-3)
}

func gRSACreate(L *lua.State) int {
    sz := L.ToInteger(1)
    pri, pub, err := RSAGenerate(sz)
    if err != nil {
        L.PushNil()
        L.PushNil()
        L.PushString(err.Error())
    } else {
        L.PushString(pri)
        L.PushString(pub)
        L.PushNil()
    }
    return 3
}

func gRSAEncrypt(L *lua.State) int {
    key := L.ToString(1)
    data := L.ToBytes(2)
    if r, err := RSAEncrypt(key, data); err != nil {
        L.PushNil()
        L.PushString(err.Error())
    } else {
        L.PushBytes(r)
        L.PushNil()
    }
    return 2
}

func gRSADecrypt(L *lua.State) int {
    key := L.ToString(1)
    data := L.ToBytes(2)
    if r, err := RSADecrypt(key, data); err != nil {
        L.PushNil()
        L.PushString(err.Error())
    } else {
        L.PushBytes(r)
        L.PushNil()
    }
    return 2
}

func gAESEncrypt(L *lua.State) int {
    aesType := L.ToString(1)
    key := L.ToBytes(2)
    data := L.ToBytes(3)
    strLen := len(key)

    if strLen == 16 || strLen == 24 || strLen == 32 {
        // ok
    } else {
        L.PushNil()
        L.PushString("key size must by 16/24/32 byte")
        return 2
    }
    switch aesType {
    case "cbc":
        if r, err := AESCBCEncrypt(key, data); err == nil {
            L.PushBytes(r)
            L.PushNil()
            return 2
        } else {
            L.PushNil()
            L.PushString(err.Error())
            return 2
        }
    case "cfb":
        if r, err := AESCFBEncrypt(key, data); err == nil {
            L.PushBytes(r)
            L.PushNil()
            return 2
        } else {
            L.PushNil()
            L.PushString(err.Error())
            return 2
        }
    case "ecb":
        if r, err := AESECBEncrypt(key, data); err == nil {
            L.PushBytes(r)
            L.PushNil()
            return 2
        } else {
            L.PushNil()
            L.PushString(err.Error())
            return 2
        }
    default:
        L.PushNil()
        L.PushString("unknown aes type, only support cbc/cfb/ecb")
        return 0
    }
}

func gAESDecrypt(L *lua.State) int {
    aesType := L.ToString(1)
    key := L.ToBytes(2)
    data := L.ToBytes(3)
    strLen := len(key)
    if strLen == 16 || strLen == 24 || strLen == 32 {
        // ok
    } else {
        L.PushNil()
        L.PushString("key size must by 16/24/32 byte")
        return 2
    }
    switch aesType {
    case "cbc":
        if r, err := AESCBCDecrypt(key, data); err == nil {
            L.PushBytes(r)
            L.PushNil()
            return 2
        } else {
            L.PushNil()
            L.PushString(err.Error())
            return 2
        }
    case "cfb":
        if r, err := AESCFBDecrypt(key, data); err == nil {
            L.PushBytes(r)
            L.PushNil()
            return 2
        } else {
            L.PushNil()
            L.PushString(err.Error())
            return 2
        }
    case "ecb":
        if r, err := AESECBDecrypt(key, data); err == nil {
            L.PushBytes(r)
            L.PushNil()
            return 2
        } else {
            L.PushNil()
            L.PushString(err.Error())
            return 2
        }
    default:
        L.PushNil()
        L.PushString("unknown aes type, only support cbc/cfb/ecb")
        return 0
    }
}

// AES
func AESCBCEncrypt(key []byte, data []byte) ([]byte, error) {
    // 分组秘钥
    // NewCipher该函数限制了输入k的长度必须为16, 24或者32
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    blockSize := block.BlockSize()                              // 获取秘钥块的长度
    data = pkcs5Padding(data, blockSize)                        // 补全码
    blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
    encrypted := make([]byte, len(data))                        // 创建数组
    blockMode.CryptBlocks(encrypted, data)                      // 加密
    return encrypted, nil
}
func AESCBCDecrypt(key []byte, data []byte) (decrypted []byte, err error) {
    var block cipher.Block
    block, err = aes.NewCipher(key)                             // 分组秘钥
    if err != nil {
        return nil, err
    }
    blockSize := block.BlockSize()                              // 获取秘钥块的长度
    blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
    decrypted = make([]byte, len(data))                         // 创建数组
    blockMode.CryptBlocks(decrypted, data)                      // 解密
    decrypted = pkcs5UnPadding(decrypted)                       // 去除补全码
    return decrypted, nil
}
func AESECBEncrypt(key []byte, data []byte) (encrypted []byte, err error) {
    var block cipher.Block
    block, err = aes.NewCipher(generateKey(key))
    if err != nil {
        return nil, err
    }
    length := (len(data) + aes.BlockSize) / aes.BlockSize
    plain := make([]byte, length*aes.BlockSize)
    copy(plain, data)
    pad := byte(len(plain) - len(data))
    for i := len(data); i < len(plain); i++ {
        plain[i] = pad
    }
    encrypted = make([]byte, len(plain))
    // 分组分块加密
    for bs, be := 0, block.BlockSize(); bs <= len(data); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
        block.Encrypt(encrypted[bs:be], plain[bs:be])
    }

    return encrypted, nil
}
func AESECBDecrypt(key []byte, data []byte) (decrypted []byte, err error) {
    var block cipher.Block
    block, err = aes.NewCipher(generateKey(key))
    if err != nil {
        return nil, err
    }
    decrypted = make([]byte, len(data))
    //
    for bs, be := 0, block.BlockSize(); bs < len(data); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
        block.Decrypt(decrypted[bs:be], data[bs:be])
    }

    trim := 0
    if len(decrypted) > 0 {
        trim = len(decrypted) - int(decrypted[len(decrypted)-1])
    }

    return decrypted[:trim], nil
}
func AESCFBEncrypt(key []byte, data []byte) (encrypted []byte, err error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    encrypted = make([]byte, aes.BlockSize+len(data))
    iv := encrypted[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }
    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(encrypted[aes.BlockSize:], data)
    return encrypted, nil
}
func AESCFBDecrypt(key []byte, data []byte) (decrypted []byte, err error) {
    block, _ := aes.NewCipher(key)
    if len(data) < aes.BlockSize {
        return nil, errors.New("ciphertext too short")
    }
    iv := data[:aes.BlockSize]
    data = data[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(data, data)
    return data, nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
    padding := blockSize - len(ciphertext)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}
func pkcs5UnPadding(origData []byte) []byte {
    length := len(origData)
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}
func generateKey(key []byte) (genKey []byte) {
    genKey = make([]byte, 16)
    copy(genKey, key)
    for i := 16; i < len(key); {
        for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
            genKey[j] ^= key[i]
        }
    }
    return genKey
}
// RSA
func RSAGenerate(bitSize int) (string, string, error) {
    // private key
    privKey, err := rsa.GenerateKey(rand.Reader, bitSize)
    if err != nil {
        return "", "", err
    }
    derStream := x509.MarshalPKCS1PrivateKey(privKey)
    block := &pem.Block{
        Type:    "RSA PRIVATE KEY",
        Headers: nil,
        Bytes:   derStream,
    }
    privBuff := bytes.NewBuffer(nil)
    err = pem.Encode(privBuff, block)
    if err != nil {
        return "", "", err
    }
    // public key
    pubKey := &privKey.PublicKey
    derPkix, err := x509.MarshalPKIXPublicKey(pubKey)
    if err != nil {
        return "", "", err
    }

    block = &pem.Block{
        Type:    "PUBLIC KEY",
        Headers: nil,
        Bytes:   derPkix,
    }
    pubBuff := bytes.NewBuffer(nil)
    err = pem.Encode(pubBuff, block)
    if err != nil {
        return "", "", err
    }
    return privBuff.String(), pubBuff.String(), nil
}
func RSAEncrypt(publicKey string, in []byte) ([]byte, error) {
    block, _ := pem.Decode([]byte(publicKey))
    if block == nil {
        return nil, errors.New("public key decode fail")
    }
    pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return nil, err
    }
    pub, ok := pubKey.(*rsa.PublicKey)
    if !ok {
        return nil, errors.New("public error")
    }
    return rsa.EncryptPKCS1v15(rand.Reader, pub, in)
}

func RSADecrypt(privateKey string, in []byte) ([]byte, error) {
    block, _ := pem.Decode([]byte(privateKey))
    if block == nil {
        return nil, errors.New("private key decode fail")
    }
    priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return nil, err
    }
    return rsa.DecryptPKCS1v15(rand.Reader, priKey, in)
}