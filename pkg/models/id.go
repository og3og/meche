package models

import (
	"crypto/rand"
	"encoding/hex"
)

func generateID() string {
	b := make([]byte, 4) // 4 bytes = 8 hex characters
	rand.Read(b)
	return hex.EncodeToString(b)
}
