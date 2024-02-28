package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func RandStr(size int) string {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
