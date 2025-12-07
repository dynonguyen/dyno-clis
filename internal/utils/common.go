package utils

import (
	"crypto/rand"
)

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	charsetLen := len(charset)
	for i := range b {
		b[i] = charset[int(b[i])%charsetLen]
	}
	return string(b)
}
