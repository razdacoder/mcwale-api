package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/razdacoder/mcwale-api/utils"
)

func CreateJWT(secret []byte, userId uuid.UUID) (string, error) {

	expiration := time.Second * time.Duration(utils.ParseStringToInt(os.Getenv("JWT_EXP"), 604800))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":     userId,
		"expires_at": time.Now().Add(expiration).Unix(),
	})

	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
