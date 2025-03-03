package activity_snap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fuzailAhmad123/test_report/infra/mongodb"
	"github.com/fuzailAhmad123/test_report/lib"
	rc "github.com/fuzailAhmad123/test_report/module/constants" //report constants
	"github.com/fuzailAhmad123/test_report/module/model"
	rsrvc "github.com/fuzailAhmad123/test_report/module/report" //report service
	rt "github.com/fuzailAhmad123/test_report/module/report/types"
	"github.com/fuzailAhmad123/test_report/module/types" //report types
	"github.com/fuzailAhmad123/test_report/module/validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ActivitySnapAggregateController(rs *types.HTTPAPIResource) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var args types.ActivitySnapAggregateArgs

		validationRes, validationErr := validation.ValidateRequestBody(r, &args)
		if validationErr != nil {
			lib.HandleError(validationRes.Message, validationRes.HttpStatus, w)
			return
		}

		//assuming - provided date by cron which is already yesterday's date of provided date
		createDate, err := lib.GetParsedTime(args.Date)
		if err != nil {
			lib.HandleError(fmt.Sprintf("Error occured while converting date"), http.StatusBadRequest, w)
			return
		}

		startDate := createDate
		endDate := startDate.Add(24 * time.Hour)

		// check if snap already present for the provided date range if present and wants to refesh  then delete them and again fetch and insert
		existsFil := bson.M{
			rc.ORG_ID: mongodb.GetOptimisticObjectIdFromHex(args.OrgId),
			rc.DATE:   bson.M{rc.GREATER_THAN_EQUALS: startDate, rc.LESSER_THAN_EQUALS: endDate},
		}

		exists, err := model.FindOne[model.ActivitySnap](context.Background(), rs.DefaultMongoDb, existsFil, options.FindOne().SetProjection(bson.M{rc.ID: 1, rc.ORG_ID: 1}))
		if err != nil {
			lib.HandleError(fmt.Sprintf("Error occured while finding snaps for particular date"), http.StatusInternalServerError, w)
			return
		}

		if exists != nil && !args.Refresh {
			response := types.ApiResponse{
				Success: true,
				Message: fmt.Sprintf("Activity Snapshot already exists for (Date: %s).", createDate.Format("2006-01-02")),
				Data:    nil,
			}

			w.WriteHeader(http.StatusOK) // Set HTTP 201 OK
			json.NewEncoder(w).Encode(response)
		} else {
			//delete the old snaps of that date
			err := model.DeleteMany[model.ActivitySnap](context.Background(), rs.DefaultMongoDb, existsFil)
			if err != nil {
				lib.HandleError(fmt.Sprintf("Error occured while deleting snaps for particular date"), http.StatusInternalServerError, w)
				return
			}
		}

		reportRes, err := rsrvc.GetReport(rsrvc.NewReportService(rs, r, true), &rt.GetActivityReportArgs{
			Start:       startDate,
			End:         endDate,
			GroupBy:     rc.ALLOWED_GROUP_BY_FOR_REPORTING,
			Metrics:     rc.ALLOWED_METRICS_FOR_REPORTING,
			OrgID:       args.OrgId,
			CampaignIds: []primitive.ObjectID{},
		})
		if err != nil {
			lib.HandleError(reportRes.Message, reportRes.HttpStatus, w)
			return
		}

		actSnapErr := CreateActivitySnapShots(rs, reportRes.Data.Report.Records, createDate)
		if actSnapErr != nil {
			lib.HandleError(fmt.Sprintf("Error occured while creating activity snapshot, %s", err.Error()), http.StatusInternalServerError, w)
			return
		}

		response := types.ApiResponse{
			Success: true,
			Message: fmt.Sprintf("Activity Snapshot created for (Date: %s) successfully", createDate.Format("2006-01-02")),
			Data:    nil,
		}

		w.WriteHeader(http.StatusCreated) // Set HTTP 201 OK
		json.NewEncoder(w).Encode(response)
	}
}
