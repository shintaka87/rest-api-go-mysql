package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"rest-api-go-mysql/db"
	"rest-api-go-mysql/models"
	"rest-api-go-mysql/utils"

	"github.com/gorilla/mux"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	result, err := db.Db.Query("SElECT * FROM users")
	if err != nil {
		log.Println(fmt.Errorf("error at Query: %w", err))
	}

	defer result.Close()

	for result.Next() {
		var user models.User

		err := result.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Println(fmt.Errorf("error at Scan: %w", err))
		}
		users = append(users, user)
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		log.Println(fmt.Errorf("error at JSONMarshal: %w", err))
	}

	utils.ReturnJsonResponse(w, http.StatusOK, usersJSON)

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.InvalidContentType)
		return
	}

	stmt, err := db.Db.Prepare("INSERT INTO users(name, email, created_at, updated_at) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Println(fmt.Errorf("error at Query: %w", err))
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Errorf("error at ReadAll: %w", err))
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	name := keyVal["name"]
	email := keyVal["email"]
	createdAt := time.Now()
	updatedAt := time.Now()

	if name == "" || email == "" {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.NameOrEmailIsEmpty)
		return
	}

	// ユーザー作成時のemailが重複していないかの確認
	if !utils.CheckEmailExist(email) {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.InvalidEmail)
		return
	}

	_, err = stmt.Exec(name, email, createdAt, updatedAt)
	if err != nil {
		log.Println(fmt.Errorf("error at stmt.Exec: %w", err))
	}

	utils.ReturnJsonResponse(w, http.StatusCreated, utils.UserCreateSuccess)

}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// パラメータで指定されたIDのuserが存在するかチェック
	if !utils.CheckUserExist(params["id"]) {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.UserNotFound)
		return
	}

	result, err := db.Db.Query("SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?", params["id"])
	if err != nil {
		log.Println(fmt.Errorf("error at Query: %w", err))
	}

	defer result.Close()
	var user models.User
	for result.Next() {
		err := result.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Println(fmt.Errorf("error at Scan: %w", err))
		}
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Println(fmt.Errorf("error at JSONMarshal: %w", err))
	}

	utils.ReturnJsonResponse(w, http.StatusOK, userJSON)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.InvalidContentType)
		return
	}

	params := mux.Vars(r)

	// パラメータで指定されたuserIDが存在するかの確認
	if !utils.CheckUserExist(params["id"]) {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.UserIdNotFound)
		return
	}

	stmt, err := db.Db.Prepare("UPDATE users SET name = ?, email = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		log.Println(fmt.Errorf("error at Query: %w", err))
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Errorf("error at ReadAll: %w", err))
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	name := keyVal["name"]
	email := keyVal["email"]
	updatedAt := time.Now()

	if name == "" && email == "" {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.NameAndEmailIsEmpty)
		return
	} else if name == "" {
		stmt, err := db.Db.Prepare("UPDATE users SET email = ?, updated_at = ? WHERE id = ?")
		if err != nil {
			log.Println(fmt.Errorf("error at Query1: %w", err))
		}
		_, err = stmt.Exec(email, updatedAt, params["id"])
		if err != nil {
			log.Println(fmt.Errorf("error at stmt.Exec: %w", err))
		}
		utils.ReturnJsonResponse(w, http.StatusOK, utils.UserUpdateSuccess)
		return
	} else if email == "" {
		stmt, err := db.Db.Prepare("UPDATE users SET name = ?, updated_at = ? WHERE id = ?")
		if err != nil {
			log.Println(fmt.Errorf("error at Query2: %w", err))
		}
		_, err = stmt.Exec(name, updatedAt, params["id"])
		if err != nil {
			log.Println(fmt.Errorf("error at stmt.Exec: %w", err))
		}
		utils.ReturnJsonResponse(w, http.StatusOK, utils.UserUpdateSuccess)
		return
	}

	_, err = stmt.Exec(name, email, updatedAt, params["id"])
	if err != nil {
		log.Println(fmt.Errorf("error at stmt.Exec: %w", err))
	}

	utils.ReturnJsonResponse(w, http.StatusAccepted, utils.UserUpdateSuccess)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// パラメータで指定されたuserIDが存在するかの確認
	if !utils.CheckUserExist(params["id"]) {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.UserIdNotFound)
		return
	}

	stmt, err := db.Db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		log.Println(fmt.Errorf("error at Query: %w", err))
	}

	_, err = stmt.Exec(params["id"])
	if err != nil {
		log.Println(fmt.Errorf("error at stmt.Exec: %w", err))
	}

	resMessage := []byte(`{
		"message": "User is deleted!"
	}`)
	utils.ReturnJsonResponse(w, http.StatusAccepted, resMessage)
}
