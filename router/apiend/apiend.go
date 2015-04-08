package apiend

import (
	"crypto/rand"
)

// Constructor for endpoint
func New() (a string) {
	a = randAPIEndPoint()
	return a
}

// Rand string generator I stole off Sourceforge.
func randAPIEndPoint() (a string) {
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, 32)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
