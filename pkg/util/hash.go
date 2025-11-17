package util

import (
	"crypto/sha256"
	"fmt"
)

func HashSHA256(data string) string {
	sum := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", sum)
}
