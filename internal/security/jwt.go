package security

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"crypto/rand"
	"encoding/base64"
)

func GenerateJWT(usuTxID uuid.UUID, email string) (string, error) {
	
	secretKey := os.Getenv("JWT_SECRET")
	
	claims := jwt.MapClaims{
		"usuTxId": usuTxID,
		"usuTxEmail":email,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
		"iat":     time.Now().Unix(),                     
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func GerarRefreshToken() (string, error) {

	refreshTokenBytes := make([] byte,32)


	_, err := rand.Read(refreshTokenBytes)

	if err != nil {
		return "", fmt.Errorf("falha ao gerar o refresh token: %w", err)
	}

	rfBase64 := base64.URLEncoding.EncodeToString(refreshTokenBytes)

	return rfBase64,nil

}