package library

import (
	"appone/models"
	"encoding/json"
	"net/http"
)

func Json(writer http.ResponseWriter, apiResponse models.ApiResponse) {

	resByte, _ := json.Marshal(apiResponse)

	writer.WriteHeader(200)
	writer.Header().Add("Content-type", "application/json")
	writer.Write(resByte)
}
