package utils

import (
	"crypto/rand"
	"crypto/sha256"
)

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// threshold is the largest multiple of len(charset) that fits in a byte (0–255).
// Bytes >= threshold are rejected to eliminate modulo bias.
// 256 / 62 = 4, so threshold = 4 * 62 = 248.
const threshold = (256 / len(charset)) * len(charset)

func GenerateHash(input string) (string, error) {
	saltBytes := make([]byte, 16)
	if _, err := rand.Read(saltBytes); err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write([]byte(input))
	h.Write(saltBytes)
	digest := h.Sum(nil) // 32 bytes

	result := make([]byte, 0, 10)
	for _, b := range digest {
		if int(b) < threshold {
			result = append(result, charset[int(b)%len(charset)])
			if len(result) == 10 {
				return string(result), nil
			}
		}
	}

	// Extremely rare: digest didn't yield 10 unbiased bytes — draw extras one at a time.
	buf := make([]byte, 1)
	for len(result) < 10 {
		if _, err := rand.Read(buf); err != nil {
			return "", err
		}
		if int(buf[0]) < threshold {
			result = append(result, charset[int(buf[0])%len(charset)])
		}
	}

	return string(result), nil
}
