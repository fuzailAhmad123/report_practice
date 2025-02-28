package lib

import (
	"encoding/json"
	"net/http"
)

func HandleError(err string, errorStatus int, w http.ResponseWriter) {
	if err != "" {
		errorRes := map[string]any{
			"success": false,
			"error":   err,
		}
		jsonResponse, jsonErr := json.Marshal(errorRes)
		if jsonErr != nil {
			http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errorStatus)
		w.Write(jsonResponse)
		return
	}
}
