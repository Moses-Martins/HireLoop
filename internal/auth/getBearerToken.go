package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
    authHeader := headers.Get("Authorization")
    if authHeader == "" {
        return "", fmt.Errorf("authorization header not found")
    }

    const prefix = "Bearer "
    authHeader = strings.TrimSpace(authHeader)

    // Check if it starts with "Bearer "
    if !strings.HasPrefix(authHeader, prefix) {
        return "", fmt.Errorf("authorization header does not contain Bearer token")
    }

    // Strip "Bearer " prefix and trim any extra spaces
    token := strings.TrimSpace(authHeader[len(prefix):])

    // Optional: check if token is empty
    if token == "" {
        return "", fmt.Errorf("authorization header contains no token")
    }

    return token, nil

}
