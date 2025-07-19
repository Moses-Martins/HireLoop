package auth

import (
	"fmt"
    "time"

    "github.com/golang-jwt/jwt/v5" 
    "github.com/google/uuid" 
)
 
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer: "hireloop",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),
    }

    // Create the token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Sign the token with the secret
    signedToken, err := token.SignedString([]byte(tokenSecret))
    if err != nil {
        return "", err
    }

    return signedToken, nil
}


func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	
    claims := &jwt.RegisteredClaims{}

    // Parse the token
    _, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        // Validate the signing algorithm
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(tokenSecret), nil
    })
 
	if err != nil {
        return uuid.Nil, err
    }

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
        return uuid.Nil, fmt.Errorf("invalid UUID in token subject: %w", err)
    }

	return userID, err
}
