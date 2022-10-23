package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"rest-api-go-mysql/db"
	"rest-api-go-mysql/models"
	"rest-api-go-mysql/utils"
)

func SearchPost(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.InvalidContentType)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Errorf("error at ReadAll: %w", err))
	}
	Keyval := make(map[string]string)
	json.Unmarshal(body, &Keyval)
	word := Keyval["word"]
	var posts []models.Post

	result, err := db.Db.Query("SELECT * FROM posts WHERE title LIKE CONCAT('%', ?, '%')", word)
	if err != nil {
		log.Println(fmt.Errorf("error at Query: %w", err))
	}

	for result.Next() {
		var post models.Post
		err := result.Scan(&post.ID, &post.Title, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			log.Println(fmt.Errorf("error at Scan: %w", err))
		}
		posts = append(posts, post)

	}

	if posts == nil {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.NoSearchResult)
		return
	}
	postsJSON, err := json.Marshal(&posts)
	if err != nil {
		log.Println(fmt.Errorf("error at JSONMarshal: %w", err))
	}

	utils.ReturnJsonResponse(w, http.StatusOK, postsJSON)

}

// func SearchPost(w http.ResponseWriter, r *http.Request) {
// 	contentTyep := r.Header.Get("Content-Type")
// 	if contentTyep != "application/json" {
// 		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.InvalidContentType)
// 		return
// 	}

// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		log.Println(fmt.Errorf("error at ReadAll: %w", err))
// 	}
// 	Keyval := make(map[string]string)
// 	json.Unmarshal(body, &Keyval)
// 	word := Keyval["word"]
// 	var posts []models.Post

// 	result, err := db.Db.Query("SELECT * FROM posts")
// 	if err != nil {
// 		log.Println(fmt.Errorf("error at Query: %w", err))
// 	}

// 	for result.Next() {
// 		var post models.Post
// 		err := result.Scan(&post.ID, &post.Title, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
// 		if err != nil {
// 			log.Println(fmt.Errorf("error at Scan: %w", err))
// 		}
// 		regex := regexp.MustCompile(word)
// 		if regex.MatchString(post.Title) {
// 			posts = append(posts, post)
// 		}
// 	}

// 	postsJSON, err := json.Marshal(&posts)
// 	if err != nil {
// 		log.Println(fmt.Errorf("error at JSONMarshal: %w", err))
// 	}

// 	utils.ReturnJsonResponse(w, http.StatusOK, postsJSON)

// }
