package shared

import (
	"encoding/json"
	"net/http"
)

func SendJSONError(w http.ResponseWriter, responseStruct any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(responseStruct)
}
