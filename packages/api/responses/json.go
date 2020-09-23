package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	Error string `json:"error"`
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		error := Error{}
		error.Error = err.Error()

		JSON(w, statusCode, error)
		return
	}
	JSON(w, http.StatusBadRequest, nil)
}
