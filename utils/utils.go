package utils

import (
	"crypto/rand"
	"fmt"
)

func NewUUID() string {
	uuid, _ := GenerateUUID()
	return uuid
}

func GenerateRandomBytes(size int) ([]byte, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return nil, fmt.Errorf("failed to read random bytes: %v", err)
	}
	return buf, nil
}

const uuidLen = 16

// GenerateUUID is used to generate a random UUID
func GenerateUUID() (string, error) {
	buf, err := GenerateRandomBytes(uuidLen)
	if err != nil {
		return "", err
	}
	return FormatUUID(buf)
}

func FormatUUID(buf []byte) (string, error) {
	if bufSize := len(buf); bufSize != uuidLen {
		return "", fmt.Errorf("wrong length byte slice (%d)", bufSize)
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		buf[0:4],
		buf[4:6],
		buf[6:8],
		buf[8:10],
		buf[10:16]), nil
}
