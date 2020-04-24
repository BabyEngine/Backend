local value = 'Hello world'

local keySize = 1024
local priv, pub, err = BabyEngine.Cipher.RSACreate(keySize)

if err ~= nil then
    print('generate rsa error:', err)
    return
end

local data, err = BabyEngine.Cipher.RSAEncrypt(pub, value)
if err ~= nil then
    print('rsa encrypt error:', err)
    return
end
local r, err = BabyEngine.Cipher.RSADecrypt(priv, data)
if err ~= nil then
    print('rsa decrypt error:', err)
    return
end

print('test ras success[', r, ']')