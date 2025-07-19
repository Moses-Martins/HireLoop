package auth

import (
	"crypto/rand"
    "encoding/hex"
    "fmt"
)

func MakeRefreshToken() (string, error) {
    // Create a 32-byte slice
    tokenBytes := make([]byte, 32)

    // Fill it with cryptographically secure random bytes
    _, err := rand.Read(tokenBytes)
    if err != nil {
        return "", fmt.Errorf("failed to generate random token: %w", err)
    }

    // Convert the random bytes to a hex string
    token := hex.EncodeToString(tokenBytes)

    return token, nil	
}