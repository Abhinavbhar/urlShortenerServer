package routes

import (
	mongoclient "Abhinavbhar/dub.sh/database"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	client := mongoclient.GetClient()
	userCollection := client.Database("dub").Collection("users")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var newUser mongoclient.User
	err = json.Unmarshal(body, &newUser)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	filter := bson.M{"username": newUser.Username}
	existing := userCollection.FindOne(context.TODO(), filter)
	if existing.Err() == nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	result, err := userCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	var secret = []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": newUser.Username,
		"id":       result.InsertedID.(primitive.ObjectID),
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		http.Error(w, "Failed to create JWT token", http.StatusInternalServerError)
		return
	}

	SetCookieHandler(w, tokenString)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
		"token":   tokenString,
	})
}

func SetCookieHandler(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     "sessionToken",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, &cookie)
}
