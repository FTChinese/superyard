package xhttp

import (
	"errors"
	"strings"
)

// ParseBearer extracts Authorization header.
// Authorization: Bearer
func ParseBearer(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("empty authorization header")
	}

	s := strings.SplitN(authHeader, " ", 2)

	bearerExists := (len(s) == 2) && (strings.ToLower(s[0]) == "bearer")

	if !bearerExists {
		return "", errors.New("bearer not found")
	}

	return s[1], nil
}
