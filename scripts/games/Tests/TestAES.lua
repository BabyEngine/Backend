
local key = '1234567890ABCDEF'
local value = 'Hello world'

local function runTest(t, key, value)
    local data, err = BabyEngine.Cipher.AESEncrypt(t, key, value)
    if err == nil then
        local decValue, err = BabyEngine.Cipher.AESDecrypt(t, key, data)
        if err == nil then
            print(t, 'test success[', decValue,']')
        else
            print(t, 'test fail decrypt:', err)
        end
    else
        print(t, 'test fail encrypt:', err)
    end
end

runTest('cbc', key, value)
runTest('cfb', key, value)
runTest('ecb', key, value)