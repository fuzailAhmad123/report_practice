package report

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/fuzailAhmad123/test_report/infra/mongodb"
	rc "github.com/fuzailAhmad123/test_report/module/constants"
	"github.com/fuzailAhmad123/test_report/module/model"
	"github.com/fuzailAhmad123/test_report/module/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportService struct {
	MongClient          *mongodb.MongoClient
	DefaultMongoDb      *mongodb.MongoDefaultDatabase
	Context             context.Context
	UseLiveActivityData bool
	UseSnapActivityData bool
	UseClickhouseData   bool
}

type GetActivityReportArgs struct {
	Start       time.Time
	End         time.Time
	GroupBy     []string
	Metrics     []string
	CampaignIds []primitive.ObjectID
	OrgID       string
}

func NewReportService(hr *types.HTTPAPIResource, r *http.Request, isInternalRequest bool) *ReportService {
	var useLiveActivityData, _ = strconv.ParseBool(r.URL.Query().Get("useLiveActivityData"))
	var useSnapActivityData, _ = strconv.ParseBool(r.URL.Query().Get("useSnapActivityData"))
	var useClickhouseData, _ = strconv.ParseBool(r.URL.Query().Get("useClickhouseData"))

	if isInternalRequest {
		useLiveActivityData = true
	}

	return &ReportService{
		MongClient:     hr.MongClient,
		DefaultMongoDb: hr.DefaultMongoDb,
		// Context:             hr.Context,
		UseLiveActivityData: useLiveActivityData,
		UseSnapActivityData: useSnapActivityData,
		UseClickhouseData:   useClickhouseData,
		// Source:              source,
	}

}

func GetReport(rs *ReportService, reportArgs *GetActivityReportArgs) (*types.ReportApiResponse, error) {
	//default resp
	response := types.ReportApiResponse{
		Success: false,
		Data:    types.ReportResponse{},
	}

	err := validate(reportArgs.Start, reportArgs.End, reportArgs.GroupBy, reportArgs.Metrics, rs)
	if err != nil {
		response.Message = err.Error()
		response.HttpStatus = http.StatusBadRequest
		return &response, err
	}

	var activityData []model.ActivityReport
	var actErr error

	if rs.UseLiveActivityData || rs.UseSnapActivityData {
		//fetch from mongo either "live" or from "snaps" based on provided type [handlin internally].
		activityData, actErr = GetActivityDataFromMongo(rs, reportArgs)
		if actErr != nil {
			return &response, err
		}
	}
	// TODO: update to fetch from clickhouse or redis also.

	//format data and total
	records, totals := GetFormattedReportResponse(activityData, reportArgs.Metrics)

	response.Message = "Report fetched succcessfully"
	response.Data = types.ReportResponse{
		GroupBy: reportArgs.GroupBy,
		Metrics: reportArgs.Metrics,
		Start:   reportArgs.Start.Format("2006-01-02"),
		End:     reportArgs.End.Format("2006-01-02"),
		Report: types.Report{
			Columns: []string{rc.ID, rc.BETS, rc.WINS, rc.AD_ID, rc.ORG_ID, rc.DATE},
			Records: records,
			Total:   totals,
		},
	}

	return &response, nil
}

func validate(start, end time.Time, groupBy []string, metrics []string, rs *ReportService) error {
	if start.IsZero() {
		return errors.New("start date is required")
	}

	if end.IsZero() {
		return errors.New("end date is required")
	}

	if len(groupBy) == 0 {
		return errors.New("groupBy is required")
	}

	if len(metrics) == 0 {
		return errors.New("metrics is required")
	}

	if !rs.UseLiveActivityData && !rs.UseSnapActivityData && !rs.UseClickhouseData {
		return errors.New("provide at least one method from [useClickhouseData, useLiveActivityData, useSnapActivityData].")
	}

	return nil
}
