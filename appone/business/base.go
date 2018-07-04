package business

import (
	"appone/library"
	"appone/models"
	"net/http"
)

type Base struct {
	status int
	msg    string
	data   interface{}
}

func (b Base) json(writer http.ResponseWriter) {
	apiResponse := models.ApiResponse{
		Status: b.status,
		Msg:    b.msg,
		Data:   b.data,
	}
	library.Json(writer, apiResponse)
}
