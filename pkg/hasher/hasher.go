package hasher

import (
	"crypto/sha256"
	"fmt"
)

type sha256Hash struct {
	salt string
}

func NewSha256Hash(salt string) *sha256Hash {
	return &sha256Hash{salt: salt}
}

func (h *sha256Hash) Hash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt)))

}
