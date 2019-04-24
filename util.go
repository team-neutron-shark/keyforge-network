package server

import (
	"math/rand"
	"time"
)

func GenerateUUID() string {
	var buffer []byte
	const choices = "abcdef0123456789"

	buffer = make([]byte, 36)

	rand.Seed(int64(time.Now().UnixNano()))

	for x := range buffer {
		if x == 8 || x == 13 || x == 18 || x == 23 {
			buffer[x] = byte('-')
		} else {
			buffer[x] = choices[rand.Intn(len(choices))]
		}
	}
	return string(buffer)
}
