package utils

import (
	"fmt"
	"log"
	"net/http"
	"rest-api-go-mysql/db"
)

func ReturnJsonResponse(w http.ResponseWriter, httpCode int, resMessage []byte) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(resMessage)
}

func CheckUserExist(userID string) bool {
	var count int
	err := db.Db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", userID).Scan(&count)
	if err != nil {
		log.Println(fmt.Errorf("error at Query: %w", err))
	}
	if count != 1 {
		return false
	}
	return true
}

func CheckPostExist(postID string) bool {
	var count int
	err := db.Db.QueryRow("SELECT COUNT(*) FROM posts WHERE id = ?", postID).Scan(&count)
	if err != nil {
		log.Println(fmt.Errorf("error at Query", err))
	}
	if count != 1 {
		return false
	}
	return true
}

func CheckEmailExist(email string) bool {
	var count int
	err := db.Db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		log.Println(fmt.Errorf("error at Query: %w", err))
	}
	if count != 0 {
		return false
	}
	return true
}
