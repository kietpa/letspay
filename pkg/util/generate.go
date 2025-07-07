package util

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/oklog/ulid/v2"
)

func GenerateReferenceId() string {
	return ulid.Make().String()
}

func GenerateRandomHex() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)

	return hex.EncodeToString(bytes)
}
