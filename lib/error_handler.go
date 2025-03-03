package lib

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/trackier/igaming-go-utils/lib/logger"
)

func HandleError(logr *logger.CustomLogger, err string, errorStatus int, w http.ResponseWriter) {
	if err != "" {
		logr.Error(context.Background(), err)
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
