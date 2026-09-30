package response

import (
	"encoding/json"
	"net/http"
)


func WriteJson(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter,statusCode int, message string){
	WriteJson(w, statusCode, map[string]string{"error": message})
}

func ReadJson(r *http.Request, target any) error {
	return json.NewDecoder(r.Body).Decode(target)
}