package utils

import "net/http"

func ReturnJsonResponse(w http.ResponseWriter, httpCode int, resMessage []byte) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(resMessage)
}

// func CheckContentType(w http.ResponseWriter, contentType string) {
// 	if contentType != "application/json"{
// 		resMessage := []byte(`{
// 			"message": "Invalid content-type: content-type should be 'application/json'"
// 		}`)
// 		ReturnJsonResponse(w, http.StatusBadRequest, resMessage)
// 		return
// 	}
// }