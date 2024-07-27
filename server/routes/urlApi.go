package routes

import (
	mongoclient "Abhinavbhar/dub.sh/database"
	"Abhinavbhar/dub.sh/middleware"
	"Abhinavbhar/dub.sh/redis"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Data struct {
	Url string `json:"url"`
}

func Url(w http.ResponseWriter, r *http.Request) {
	//getting the mongoDb client both user and url collection
	Mongoclient := mongoclient.GetClient()
	urlCollection := Mongoclient.Database("dub").Collection("url")
	userCollection := Mongoclient.Database("dub").Collection("users")

	var id string
	//getting the id which is set by auth middleware
	if User, ok := r.Context().Value(middleware.UserKey).(middleware.User); ok {
		id = User.ID
	} else {
		http.Error(w, "login again missing credentials", http.StatusBadRequest)
		return
	}
	//reading the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	// parsing the from json data
	var data Data
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}
	//function to generate a random string
	uri := generateRandomString(4)
	finalUrl := "http://localhost:8080/" + uri
	//converting id from string to mongo primitive
	userID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Fatal(err)
	}
	//making the new url to store it in mongoDb
	var url mongoclient.ActiveLink
	url.URL = data.Url
	url.ShortCode = uri
	url.UserId = userID
	update := bson.M{
		"$push": bson.M{"active_links": url},
	}
	//finding the user from id
	var user mongoclient.User
	filter := bson.M{
		"_id": userID}
	error12 := userCollection.FindOne(context.TODO(), filter).Decode(&user)
	if error12 != nil {
		fmt.Println(error12)
		http.Error(w, "user does not exist ", http.StatusBadRequest)
		return
	}
	//checking if url already exists in the user
	for _, link := range user.ActiveLinks {
		if link.URL == data.Url {
			http.Error(w, "url already exists", http.StatusBadRequest)
			return
		}
	}
	//if links are less then 10
	if len(user.ActiveLinks) >= 10 {
		http.Error(w, "url limit exceed only 10 url are accepted", http.StatusBadRequest)
		return
	}
	//everything is ok store in redis and mongoDb
	client := redis.RedisDatabase()
	ctx := r.Context()
	if err := client.Set(ctx, uri, data.Url, 1*time.Hour).Err(); err != nil {
		http.Error(w, "Failed to store URL in Redis", http.StatusInternalServerError)
		return
	}
	_, error := urlCollection.InsertOne(context.TODO(), url)
	_, error1 := userCollection.UpdateByID(context.TODO(), userID, update)

	if error != nil {
		http.Error(w, "failed to store your url", http.StatusInternalServerError)
		return
	}
	if error1 != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	//make the response and send the final response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"finalUrl": finalUrl, "id": id}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	charset := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
