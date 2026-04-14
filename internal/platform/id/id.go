package id

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func New(prefix string) string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%s-fallback", prefix)
	}
	return prefix + "-" + hex.EncodeToString(b)
}
