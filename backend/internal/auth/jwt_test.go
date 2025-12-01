package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	userID := uint(1)
	role := "admin"

	tokenString, err := GenerateToken(userID, role)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &Claims{})
	assert.NoError(t, err)
	assert.NotNil(t, token)
}

func TestValidateToken_Valid(t *testing.T) {
	userID := uint(1)
	role := "admin"

	tokenString, _ := GenerateToken(userID, role)

	claims, err := ValidateToken(tokenString)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, role, claims.Role)
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	userID := uint(1)
	role := "admin"

	tokenString, _ := GenerateToken(userID, role)

	invalidToken := tokenString[:len(tokenString)-1] + "X"

	claims, err := ValidateToken(invalidToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_Expired(t *testing.T) {
	expirationTime := time.Now().Add(-1 * time.Hour) // Expired 1 hour ago
	claims := &Claims{
		UserID: 1,
		Role:   "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)

	validatedClaims, err := ValidateToken(tokenString)

	assert.Error(t, err)
	assert.Nil(t, validatedClaims)
	assert.Contains(t, err.Error(), "token is expired")
}
