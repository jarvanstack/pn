package utils

import (
	"math/rand"
	"time"
)

var (
	letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func init() {
	rand.Seed(time.Now().UnixNano())

}

func RandString(size int) string {
	return string(RandBytes(size))
}

func RandInt(max int) int {
	return rand.Intn(max)
}

func RandBytes(size int) []byte {
	b := make([]byte, size)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}
