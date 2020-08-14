package util

import (
	"crypto/md5"
	"encoding/hex"
)

func StringHash(s string) (hash uint16) {
	for _, c := range s {
		ch := uint16(c)
		hash = hash + ((hash) << 5) + ch + (ch << 7)
	}

	return
}

func BytesMD5(data []byte) string {
	m := md5.New()
	m.Write(data)

	return hex.EncodeToString(m.Sum(nil))
}

func StringMD5(s string) string {
	return BytesMD5([]byte(s))
}