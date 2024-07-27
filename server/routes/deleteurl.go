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

type deleteUrl struct {
	URL string `json:"short_code" bson:"short_code"`
}

func DeleteUrl(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user data from the context
	data, ok := r.Context().Value(middleware.UserKey).(middleware.User)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Println("User data:", data)

	// Convert the user ID to an ObjectID
	objectId, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	// Get the MongoDB client and collection
	client := mongoClient.GetClient()
	userCollection := client.Database("dub").Collection("users")
	urlCollection := client.Database("dub").Collection("url")

	// Read and parse the request body
	body, parsingError := io.ReadAll(r.Body)
	if parsingError != nil {
		http.Error(w, "error parsing request", http.StatusBadRequest)
		return
	}
	var url deleteUrl
	err = json.Unmarshal(body, &url)
	if err != nil {
		http.Error(w, "error unmarshalling request body", http.StatusBadRequest)
		return
	}
	fmt.Println("URL to delete:", url.URL)

	// Create the filter and update documents
	filter := bson.M{"_id": objectId}
	fmt.Println(url.URL)
	update := bson.M{"$pull": bson.M{
		"active_links": bson.M{
			"short_code": url.URL,
		},
	}}

	// Perform the update operation
	updateResult, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	urlFilter := bson.M{"short_code": url.URL}
	_, error := urlCollection.DeleteOne(context.TODO(), urlFilter)
	if error != nil {
		http.Error(w, "internsal server error", http.StatusInternalServerError)
		return
	}

	// Check if any documents were modified
	if updateResult.ModifiedCount == 0 {
		http.Error(w, "no document found or URL not present", http.StatusNotFound)
		return
	}

	// Send success response
	response := bson.M{
		"message": "delete successful",
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "error marshalling response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}
