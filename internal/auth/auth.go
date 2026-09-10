package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("Error hashing password: %s", err)
	}

	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("Error comparing password to hash: %s", err)
	}

	return match, nil
}

type chirpyAccessClaims struct {
	
	jwt.RegisteredClaims
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := chirpyAccessClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Issuer:    "chirpy-access",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	if ss == "" {
		return "", fmt.Errorf("Could not create signed string")
	}

	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &chirpyAccessClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(subject)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	token, ok := strings.CutPrefix(authHeader, "Bearer ")
	if !ok {
		return "", fmt.Errorf("Bearer token not found")
	}
	
	return token, nil
}

func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)
	
	return hex.EncodeToString(key)
}