package routes

import (
	"Abhinavbhar/dub.sh/middleware"
	"encoding/json"
	"net/http"
)

func AuthChecker(w http.ResponseWriter, r *http.Request) {
	var respone map[string]string

	if u, ok := r.Context().Value(middleware.UserKey).(middleware.User); ok {
		respone = map[string]string{"username": u.Username, "id": u.ID}
	} else {
		http.Error(w, "please login again, username not found in the request", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(respone)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	w.Write(jsonResponse)
}
