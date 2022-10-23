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
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.InvalidContentType)
		return
	}

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

	if title == "" || userID == 0 {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.TitleOrIdIsEmpty)
		return
	}

	// POSTで送られてきたuser_idのuserが存在するかチェック
	if !utils.CheckUserExist(keyVal["user_id"]) {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.UserIdNotFound)
		return
	}

	_, err = stmt.Exec(title, userID, createdAt, updatedAt)
	if err != nil {
		log.Println(fmt.Errorf("error at stmt.Exec: %w", err))
	}

	utils.ReturnJsonResponse(w, http.StatusCreated, utils.PostCreateSuccess)

}

func GetPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	fmt.Printf("%T\n", params)

	// パラメータに指定されたpostのIDが存在するかのチェック
	if !utils.CheckPostExist(params["id"]) {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.PostNotFound)
		return
	}

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
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.InvalidContentType)
		return
	}

	params := mux.Vars(r)

	// パラメータに指定されたIDのpostが存在するかのチェック
	if !utils.CheckPostExist(params["id"]) {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.PostNotFound)
		return
	}

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
	if newTitle == "" {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.TitleIsEmpty)
		return
	}
	updatedAt := time.Now()

	_, err = stmt.Exec(newTitle, updatedAt, params["id"])
	if err != nil {
		log.Println(fmt.Errorf("error at stmt.Exec: %w", err))
	}

	utils.ReturnJsonResponse(w, http.StatusAccepted, utils.PostUpdateSuccess)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if !utils.CheckPostExist(params["id"]) {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.PostNotFound)
		return
	}

	stmt, err := db.Db.Prepare("DELETE FROM posts WHERE id = ?")
	if err != nil {
		log.Println(fmt.Errorf("error at Query: %w", err))
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		log.Println(fmt.Errorf("error at stmt.Exec: %w", err))
	}

	utils.ReturnJsonResponse(w, http.StatusOK, utils.PostDeleteSuccess)
}
