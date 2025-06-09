package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	data := make([]byte, 32)
	rand.Read(data)
	token := hex.EncodeToString(data)

	return token
}
