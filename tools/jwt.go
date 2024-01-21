package tools

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type GenerateJWTTokenParams struct {
	FullName string    `json:"full_name"`
	GUID     uuid.UUID `json:"user_guid"`
}

type JWTCustomClaims struct {
	GenerateJWTTokenParams
	jwt.RegisteredClaims
}

func generateToken(claim JWTCustomClaims, lifetime int, key string) (tokenString string, expiredAt time.Time, err error) {
	now := time.Now().UTC()
	timeExpiredAt := now.Add(time.Duration(lifetime) * time.Hour)

	claim.IssuedAt = jwt.NewNumericDate(now)
	claim.ExpiresAt = jwt.NewNumericDate(timeExpiredAt)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
	if err != nil {
		return
	}

	tokenString, err = token.SignedString(parsedKey)
	if err != nil {
		return
	}

	return tokenString, timeExpiredAt, nil
}

func GenerateJWTToken(params GenerateJWTTokenParams, lifetime int, key string) (token string, expiredAt time.Time, err error) {
	claim := JWTCustomClaims{GenerateJWTTokenParams: params}

	return generateToken(claim, lifetime, key)
}
