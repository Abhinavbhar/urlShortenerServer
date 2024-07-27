package main

import (
	mongoClient "Abhinavbhar/dub.sh/database"
	"Abhinavbhar/dub.sh/middleware"
	"Abhinavbhar/dub.sh/redis"
	"Abhinavbhar/dub.sh/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Configure CORS settings

	corsHandler := cors.New(cors.Options{
		//add frontend url here
		AllowedOrigins: []string{"http://localhost:3000"},
		//options
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		//allowed headers
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding",
			"X-CSRF-Token", "Authorization"},
		ExposedHeaders: []string{"Link"},
		//cookies
		AllowCredentials: true,
		MaxAge:           300,
	})

	// Initialize MongoDB and Redis
	mongoClient.InitMongoClient()
	redis.InitRedis()
	// Create a new router
	r := mux.NewRouter()
	r.Handle("/api/createurl", middleware.AuthMiddleware(http.HandlerFunc(routes.Url))).Methods("POST")
	r.Handle("/url/ip", middleware.AuthMiddleware(http.HandlerFunc(routes.IpAddress))).Methods("POST")
	r.Handle("/customurl", middleware.AuthMiddleware(http.HandlerFunc(routes.CustomUrl))).Methods("POST")
	r.Handle("/deleteurl", middleware.AuthMiddleware(http.HandlerFunc(routes.DeleteUrl))).Methods("POST")
	r.Handle("/authcheck", middleware.AuthMiddleware(http.HandlerFunc(routes.AuthChecker))).Methods("GET")
	r.Handle("/dashboard", middleware.AuthMiddleware(http.HandlerFunc(routes.Dashboard))).Methods("GET")
	r.HandleFunc("/{value}", routes.RedirectUrl).Methods("GET")
	r.HandleFunc("/login", routes.Login).Methods("POST")
	r.HandleFunc("/signup", routes.Signup).Methods("POST")
	handler := corsHandler.Handler(r)
	log.Println("starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
