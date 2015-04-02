package apiend

import (
	"crypto/rand"
)

// APIEndPoint is just a random string for endpoints.  Everytime Mongolar reboots a new endpoint is created
type APIEndPoint string

// Constructor for endpoint
func New() *APIEndPoint {
	a = new(APIEndPoint)
	a.randAPIEndPoint()
	return a
}

// Rand string generator I stole off Sourceforge.
func (a *APIEndPoint) randAPIEndPoint() {
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, 32)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	a = bytes
}
