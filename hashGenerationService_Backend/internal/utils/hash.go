package utils

import (
	"crypto/rand"
	"crypto/sha256"
)

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateHash(input string) (string, error) {
	randBytes := make([]byte, 16)
	if _, err := rand.Read(randBytes); err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write([]byte(input))
	h.Write(randBytes)
	sum := h.Sum(nil)

	result := make([]byte, 10)
	for i := 0; i < 10; i++ {
		result[i] = charset[int(sum[i])%len(charset)]
	}
	return string(result), nil
}
