package shared

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func getCode(code int) string {
	switch code {
	case http.StatusInternalServerError:
		return INTERNAL_ERROR_CODE
	case http.StatusUnauthorized:
		return UNAUTHORIZED_CODE
	case http.StatusConflict:
		return CONFLICT_CODE
	case http.StatusNotFound:
		return NOT_FOUND_CODE
	case http.StatusBadRequest:
		return BAD_REQUEST_CODE
	default:
		return INTERNAL_ERROR_CODE
	}
}

func SendJSONError(w http.ResponseWriter, responseStruct ErrorResponse, code int) {
	if responseStruct.Code == "" || len(responseStruct.Code) == 0 {
		responseStruct.Code = getCode(code)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(responseStruct)
}
