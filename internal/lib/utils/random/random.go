package random

import (
	"math/rand"
	"time"
)

func GetRandomString(length int) string {
	idx := rand.New(rand.NewSource(time.Now().UnixNano()))
	alphabet := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, length)

	for i := range b {
		b[i] = alphabet[idx.Intn(len(alphabet))]
	}

	return string(b)
}
