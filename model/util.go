package model

import (
	"math/rand"
	"time"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	idLen   = 50
)

func getNewID() string {
	b := make([]byte, idLen)
	randSeed := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range b {
		b[i] = charset[randSeed.Intn(len(charset))]
	}

	return string(b)
}
