package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"rest-api-go-mysql/db"
	"rest-api-go-mysql/models"
	"rest-api-go-mysql/utils"

	"github.com/gorilla/mux"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post

	result, err := db.Db.Query("SELECT * from posts")
	if err != nil {
		log.Println(fmt.Errorf("error at Query: %w", err))
	}

	defer result.Close()

	for result.Next() {
		var post models.Post

		err := result.Scan(&post.ID, &post.Title, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			log.Println(fmt.Errorf("error at Scan: %w", err))
		}
		posts = append(posts, post)
	}

	postsJSON, err := json.Marshal(&posts)
	if err != nil {
		log.Println(fmt.Errorf("error at JSONMarshal: %w", err))
	}

	utils.ReturnJsonResponse(w, http.StatusOK, postsJSON)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Db.Prepare("INSERT INTO posts(title, user_id, created_at, updated_at) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Println(fmt.Errorf("error at Query: %w", err))
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Errorf("error at ReadAll: %w", err))
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	title := keyVal["title"]
	userID, _ := strconv.Atoi(keyVal["user_id"])
	createdAt := time.Now()
	updatedAt := time.Now()

	_, err = stmt.Exec(title, userID, createdAt, updatedAt)
	if err != nil {
		log.Println(fmt.Errorf("error at stmt.Exec: %w", err))
	}

	resMessage := []byte(`{
		"message": "Post is created!"
	}
	`)

	utils.ReturnJsonResponse(w, http.StatusCreated, resMessage)

}

func GetPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	result, err := db.Db.Query("SELECT id, title, user_id, created_at, updated_at FROM posts WHERE id = ?", params["id"])
	if err != nil {
		log.Println(fmt.Errorf("error at Query: %w", err))
	}
	defer result.Close()

	var post models.Post

	for result.Next() {
		err := result.Scan(&post.ID, &post.Title, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			log.Println(fmt.Errorf("error at Scan: %w", err))
		}

	}
	
	postJSON, err := json.Marshal(post)
	if err != nil {
		log.Println(fmt.Errorf("error at JSONMarshal: %w", err))
	}
	utils.ReturnJsonResponse(w, http.StatusOK, postJSON)

}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	stmt, err := db.Db.Prepare("UPDATE posts SET title = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		log.Println(fmt.Errorf("error at Query: %w", err))
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Errorf("error at ReadAll: %w", err))
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newTitle := keyVal["title"]
	updatedAt := time.Now()

	_, err = stmt.Exec(newTitle, updatedAt, params["id"])
	if err != nil {
		log.Println(fmt.Errorf("error at stmt.Exec: %w", err))
	}

	resMessage := []byte(`{
		"message": "Post is updated!"
	}
	`)

	utils.ReturnJsonResponse(w, http.StatusAccepted, resMessage)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	stmt, err := db.Db.Prepare("DELETE FROM posts WHERE id = ?")
	if err != nil {
		log.Println(fmt.Errorf("error at Query: %w", err))
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		log.Println(fmt.Errorf("error at stmt.Exec: %w", err))
	}

	resMessage := []byte(`{
		"message": "Post is deleted!"
	}`)

	utils.ReturnJsonResponse(w, http.StatusOK, resMessage)
}
