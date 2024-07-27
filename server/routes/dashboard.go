package routes

import (
	mongoClient "Abhinavbhar/dub.sh/database"
	"Abhinavbhar/dub.sh/middleware"
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DashboardResponse struct {
	Username    string                   `json:"username"`
	ID          string                   `json:"id"`
	ActiveLinks []mongoClient.ActiveLink `json:"active_links"`
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	if u, ok := r.Context().Value(middleware.UserKey).(middleware.User); ok {
		client := mongoClient.GetClient()
		userCollection := client.Database("dub").Collection("users")
		objectId, err := primitive.ObjectIDFromHex(u.ID)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		filter := bson.M{"_id": objectId}
		var result mongoClient.User
		err = userCollection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
			return
		}

		response := DashboardResponse{
			Username:    u.Username,
			ID:          u.ID,
			ActiveLinks: result.ActiveLinks,
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, _ := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		w.Write(jsonResponse)
	} else {
		http.Error(w, "Please login again, username not found in the request", http.StatusInternalServerError)
	}
}
