package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})

	signedJWT, err := jwtToken.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signedJWT, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return [16]byte{}, fmt.Errorf("invalid token: %v", err)
	}

	if !token.Valid {
		return [16]byte{}, fmt.Errorf("token is not valid")
	}

	stringID, err := token.Claims.GetSubject()
	if err != nil {
		return [16]byte{}, fmt.Errorf("can`t get userID: %v", err)
	}
	userID, err := uuid.Parse(stringID)
	if err != nil {
		return [16]byte{}, fmt.Errorf("can`t parse string to UUID: %v", err)
	}
	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("token is not exist")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("token is not valid")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	return token, nil
}
