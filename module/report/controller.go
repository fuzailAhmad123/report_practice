package report

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fuzailAhmad123/test_report/lib"

	"github.com/fuzailAhmad123/test_report/module/types"
	"github.com/fuzailAhmad123/test_report/module/validation"
)

func GetActivityReportController(rs *types.HTTPAPIResource) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var args types.ReportRequestArgs

		validationRes, validationErr := validation.ValidateGetQueryParams(r, &args)
		if validationErr != nil {
			lib.HandleError(validationRes.Message, validationRes.HttpStatus, w)
			return
		}

		reportArgs, err := ConvertReportQueryParams(&args)
		if err != nil {
			lib.HandleError(fmt.Sprintf("Error while converting the args: %v", err), http.StatusBadRequest, w)
			return
		}

		reportRes, err := GetReport(NewReportService(rs, r, false), reportArgs)
		if err != nil {
			lib.HandleError(reportRes.Message, reportRes.HttpStatus, w)
			return
		}

		w.WriteHeader(http.StatusOK) // Set HTTP 200 OK
		json.NewEncoder(w).Encode(reportRes)
	}
}
