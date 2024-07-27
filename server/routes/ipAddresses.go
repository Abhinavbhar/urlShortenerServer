package routes

import (
	mongoClient "Abhinavbhar/dub.sh/database"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

type req struct {
	Short_code string `json:"short_code"`
}

func IpAddress(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read request", http.StatusBadRequest)
		return
	}
	var url req
	jsonerror := json.Unmarshal(body, &url)
	if jsonerror != nil {
		http.Error(w, "error parsing request", http.StatusBadRequest)
		return
	}
	client := mongoClient.GetClient()
	urlCollection := client.Database("dub").Collection("url")
	filter := bson.M{
		"short_code": url.Short_code,
	}
	var urlFound mongoClient.ActiveLink
	fmt.Println(url)
	notFound := urlCollection.FindOne(context.TODO(), filter).Decode(&urlFound)
	if notFound != nil {
		http.Error(w, "cannot find the link", http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(urlFound.Ip)
	w.Write(response)

}
