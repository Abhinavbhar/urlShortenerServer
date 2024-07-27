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

func Login(w http.ResponseWriter, r *http.Request) {

	client := mongoclient.GetClient()
	userCollection := client.Database("dub").Collection("users")
	body, _ := io.ReadAll(r.Body)
	var b mongoclient.User
	err := json.Unmarshal(body, &b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	filter := bson.M{"username": b.Username}
	existing := userCollection.FindOne(context.TODO(), filter)

	var result bson.M
	existenceError := existing.Decode(&result)
	id, _ := result["_id"].(primitive.ObjectID)
	username := result["username"]
	if existenceError != nil {
		http.Error(w, "username does not exist", http.StatusConflict)
		return
	}

	var secret = []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"id":       id,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		panic(err)
	}
	SetCookieHandler(w, tokenString)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "logged in successfully"})
}
