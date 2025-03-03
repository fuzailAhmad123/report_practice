package acitvity

import (
	"encoding/json"
	"net/http"

	"github.com/fuzailAhmad123/test_report/lib"
	"github.com/fuzailAhmad123/test_report/module/types"
	"github.com/fuzailAhmad123/test_report/module/validation"
)

// CreateActivityController is controller to handle create activity requests.
func CreateActivityController(rs *types.HTTPAPIResource) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		//fetch -> bets, wins, adid, orgid, date
		var args CreateActivityArgs

		validationRes, validationErr := validation.ValidateRequestBody(r, &args)
		if validationErr != nil {
			lib.HandleError(rs.Logr, validationRes.Message, validationRes.HttpStatus, w)
			return
		}

		activityRes, activityErr := CreateActivityService(rs, &args)
		if activityErr != nil {
			lib.HandleError(rs.Logr, activityRes.Message, activityRes.HttpStatus, w)
			return
		}

		w.WriteHeader(http.StatusCreated) // Set HTTP 201 CREATED
		json.NewEncoder(w).Encode(activityRes)
	}
}
