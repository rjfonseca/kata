package scaffold

import (
	"crypto/sha256"
	"fmt"
	"os"
)

// ChecksumFile calculates a SHA-256 checksum of a file.
func ChecksumFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(data)
	return fmt.Sprintf("sha256:%x", sum), nil
}
