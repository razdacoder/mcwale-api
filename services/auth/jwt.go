package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/razdacoder/mcwale-api/models"
	"github.com/razdacoder/mcwale-api/utils"
)

func CreateJWT(secret []byte, userID uuid.UUID, role string) (string, error) {
	expiration := time.Second * time.Duration(utils.ParseStringToInt(os.Getenv("JWT_EXP"), 604800))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":     userID,
		"role":       role,
		"expires_at": time.Now().Add(expiration).Unix(),
	})

	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func CreateResetPasswordJWT(secret []byte, userID uuid.UUID) (string, error) {
	expiration := time.Duration(10) * time.Minute
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":     userID,
		"expires_at": time.Now().Add(expiration).Unix(),
	})
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
func verifyToken(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return "", "", err
	}

	if !token.Valid {
		return "", "", fmt.Errorf("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)
	userId, ok := claims["userID"].(string)
	if !ok {
		return "", "", fmt.Errorf("user claim not found")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", "", fmt.Errorf("role claim not found")
	}

	return userId, role, nil

}

func VerifyPasswordToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)
	userId, ok := claims["userID"].(string)
	if !ok {
		return "", fmt.Errorf("userId claim not found")
	}

	return userId, nil

}

type ContextKey string

const UserKey ContextKey = "user"

func IsLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		tokenString := request.Header.Get("Authorization")
		if tokenString == "" {
			utils.WriteError(writer, http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
			return
		}
		tokenString = tokenString[len("Bearer "):]
		userId, role, err := verifyToken(tokenString)
		if err != nil {
			utils.WriteError(writer, http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
			return
		}
		ctxVal := map[string]string{"userId": userId, "role": role}
		ctx := context.WithValue(request.Context(), UserKey, ctxVal)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		ctx := request.Context()
		user, ok := ctx.Value(UserKey).(map[string]string)
		if !ok {
			utils.WriteError(writer, http.StatusUnprocessableEntity, fmt.Errorf("no user found"))
			return
		}

		if models.UserRole(user["role"]) != models.Admin {
			utils.WriteError(writer, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func IsAdminOrCurrentUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		ctx := request.Context()
		user, ok := ctx.Value(UserKey).(map[string]string)
		if !ok {
			utils.WriteError(writer, http.StatusUnprocessableEntity, fmt.Errorf("no user found"))
			return
		}

		id := chi.URLParam(request, "id")
		fmt.Println(user["userId"], id)
		if models.UserRole(user["role"]) != models.Admin && user["userId"] != id {
			utils.WriteError(writer, http.StatusUnauthorized, fmt.Errorf("unauthorized this"))
			return
		}
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
