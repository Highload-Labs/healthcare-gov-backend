package shared

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

func Hash(token string) (string, error) {
	if token == "" {
		return "", errors.New("cannot hash an empty token")
	}

	hasher := sha256.New()
	hasher.Write([]byte(token))
	hashBytes := hasher.Sum(nil)

	return hex.EncodeToString(hashBytes), nil
}
