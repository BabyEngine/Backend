local value = 'Hello world'

local ok, _md5, _sha1, _sha256, _sha512 = BabyEngine.Hash.HashString(value, 'md5', 'sha1', 'sha256', 'sha512')

print(ok, _md5, _sha1, _sha256, _sha512)
