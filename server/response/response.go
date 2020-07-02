package response

import (
	"encoding/json"
	"log"
	"net/http"
)

func Success(w http.ResponseWriter, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("faield to marshal json: %v", err)
		InternalServerError(w, err)
		return
	}
	w.Write(jsonData)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func BadRequest(w http.ResponseWriter, message interface{}) {
	httpError(w, http.StatusBadRequest, message)
}

func InternalServerError(w http.ResponseWriter, message interface{}) {
	httpError(w, http.StatusInternalServerError, message)
}

func httpError(w http.ResponseWriter, code int, message interface{}) {
	jsonData, _ := json.Marshal(httpErrorResponse{Code: code, Message: message})

	w.WriteHeader(code)
	w.Write(jsonData)
}

type httpErrorResponse struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}
