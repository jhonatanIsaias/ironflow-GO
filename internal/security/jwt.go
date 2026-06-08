package security

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJWT(usuTxID uuid.UUID, email string) (string, error) {
	
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "ironflow-super-secret-dev-key"
	}

	claims := jwt.MapClaims{
		"usuTxId": usuTxID,
		"usuTxEmail":email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
		"iat":     time.Now().Unix(),                     
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}