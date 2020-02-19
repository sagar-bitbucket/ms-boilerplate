package responses

import (
	"encoding/json"
	"io"
	"net/http"
)

//jsonResponse Type
type jsonResponse struct {
	Data interface{} `json:"data"`
}

//jsonErrorResponse Type
type jsonErrorResponse struct {
	Error string `json:"message"`
}

//jsonResponse Type
type jsonAllResponse struct {
	Count int         `json:"count"`
	Data  interface{} `json:"data"`
}

//ReadInput from the body
func ReadInput(rBody io.ReadCloser, input interface{}) error {
	decoder := json.NewDecoder(rBody)
	err := decoder.Decode(input)
	return err
}

//WriteOKResponse as a standard JSON response with StatusOK
func WriteOKResponse(w http.ResponseWriter, m interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&jsonResponse{Data: m}); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Internal Sever Error!")
	}
}

//WriteErrorResponse as a Standard API JSON response with a response code and error
func WriteErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)
	json.
		NewEncoder(w).Encode(&jsonErrorResponse{Error: errorMsg})
}
