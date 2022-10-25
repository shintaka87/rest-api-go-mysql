package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rest-api-go-mysql/db"
	"rest-api-go-mysql/models"
	"rest-api-go-mysql/utils"
)

func SearchPost(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	word := params.Get("word")
	if word == "" {
		return
	}
	var posts []models.Post

	result, err := db.Db.Query("SELECT * FROM posts WHERE title LIKE CONCAT('%', ?, '%')", word)
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
