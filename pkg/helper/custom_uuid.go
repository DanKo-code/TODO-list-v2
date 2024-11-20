package helper

import (
	"crypto/rand"
	"fmt"
	"io"
)

func GenerateUUID() (string, error) {
	uuid := make([]byte, 16)

	if _, err := io.ReadFull(rand.Reader, uuid); err != nil {
		return "", err
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:16]), nil
}
