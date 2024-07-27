package routes

import (
	mongoClient "Abhinavbhar/dub.sh/database"
	"Abhinavbhar/dub.sh/middleware"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type customUrl struct {
	Custom_url string `json:"custom_url"`
	Url        string `json:"url"`
}

func CustomUrl(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading request", http.StatusBadRequest)
		return
	}
	var reqUrl customUrl
	parsingError := json.Unmarshal(body, &reqUrl)
	if parsingError != nil {
		http.Error(w, "error parsing request", http.StatusBadRequest)
		return
	}

	mongo := mongoClient.GetClient()
	urlCollection := mongo.Database("dub").Collection("url")
	userCollection := mongo.Database("dub").Collection("users")
	filter := bson.M{
		"short_code": reqUrl.Custom_url,
	}
	existingError := urlCollection.FindOne(context.TODO(), filter)
	var existing bson.M
	existingErr := existingError.Decode(&existing)
	if existingErr == nil {
		http.Error(w, "sorry url is already taken", http.StatusBadRequest)
		fmt.Println(existingErr)
		return
	}
	data := r.Context().Value(middleware.UserKey).(middleware.User)
	id, _ := primitive.ObjectIDFromHex(data.ID)

	var newUrl mongoClient.ActiveLink
	fmt.Println(reqUrl)
	newUrl.UserId = id
	newUrl.ShortCode = reqUrl.Custom_url
	newUrl.URL = reqUrl.Url
	_, mongoErr := urlCollection.InsertOne(context.TODO(), newUrl)
	if mongoErr != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	userFilter := bson.M{
		"_id": id,
	}
	update := bson.M{
		"$push": bson.M{
			"active_links": newUrl,
		},
	}
	_, usererr := userCollection.UpdateOne(context.TODO(), userFilter, update)
	if usererr != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "added successfully",
	})

}
