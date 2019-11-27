package helper

import (
	"math/rand"
	"time"
)

const (
	// CHARS for setting short random string
	CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// ErrorDataNotFound error message when data doesn't exist
	ErrorDataNotFound = "data %s not found"
)

// RandomString function for random string
func RandomString(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	charsLength := len(CHARS)
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = CHARS[rand.Intn(charsLength)]
	}
	return string(result)
}
