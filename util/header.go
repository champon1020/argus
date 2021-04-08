package util

import (
	"errors"
	"strings"
)

// ExtractBearerToken extracts the token from authentication header.
func ExtractBearerToken(auth string) (string, error) {
	el := strings.Split(auth, "Bearer ")
	if len(el) < 2 {
		return "", errors.New("invalid authorization header")
	}
	return el[1], nil
}
