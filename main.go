package main

import (
	"log"
	"net/http"
	"os"

	"rest-api-go-mysql/db"
	"rest-api-go-mysql/handler"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)


var err error

func main() {

	defer db.Db.Close()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()
	router.HandleFunc("/posts", handler.GetPosts).Methods("GET")
	router.HandleFunc("/posts", handler.CreatePost).Methods("POST")
	router.HandleFunc("/posts/{id}", handler.GetPost).Methods("GET")
	router.HandleFunc("/posts/{id}", handler.UpdatePost).Methods("PUT")
	router.HandleFunc("/posts/{id}", handler.DeletePost).Methods("DELETE")
	router.HandleFunc("/users", handler.GetUsers).Methods("GET")
	router.HandleFunc("/users", handler.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")

	http.ListenAndServe(os.Getenv("GO_PORT"), router)
}
