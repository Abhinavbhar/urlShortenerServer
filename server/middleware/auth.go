package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type contextKey string

const UserKey contextKey = "user"

type User struct {
	Username string
	ID       string
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("sessionToken")
		if err != nil {
			http.Error(w, "you have been logged out please login again", http.StatusBadRequest)
			return
		}
		var value string = cookie.Value
		username, id, error := verifyToken(value)
		if error != nil {
			http.Error(w, "incorrect creentials", http.StatusBadRequest)
			return
		}
		user := User{
			Username: username,
			ID:       id,
		}
		ctx := context.WithValue(r.Context(), UserKey, user)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

	})
}

func verifyToken(tokenString string) (string, string, error) {
	godotenv.Load()
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return "", "", err
	}
	if !token.Valid {
		return "", "", fmt.Errorf("invalid token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, usernameOk := claims["username"].(string)
		id, idOk := claims["id"].(string)
		if usernameOk && idOk {
			return username, id, nil
		}
		return "", "", fmt.Errorf("invalid claims")
	}

	return "", "", fmt.Errorf("username claims not found in the token")
}
