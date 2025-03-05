package report

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	rc "github.com/fuzailAhmad123/test_report/module/constants" //report constants
	"github.com/fuzailAhmad123/test_report/module/report/retriever/clickhouse"
	mongolive "github.com/fuzailAhmad123/test_report/module/report/retriever/mongo_live"
	mongosnap "github.com/fuzailAhmad123/test_report/module/report/retriever/mongo_snap"
	"github.com/fuzailAhmad123/test_report/module/report/retriever/redis"
	rt "github.com/fuzailAhmad123/test_report/module/report/types" //report types
	"github.com/fuzailAhmad123/test_report/module/types"
)

func NewReportService(hr *types.HTTPAPIResource, r *http.Request, isInternalRequest bool) *rt.ReportService {
	var source = r.URL.Query().Get("source")

	var rr *rt.ReportRetriever

	if isInternalRequest {
		rr = &rt.ReportRetriever{Name: rc.MONGO_LIVE, Retriever: mongolive.Init()}
	} else if source != "" {
		switch source {
		case rc.CLICKHOUSE:
			rr = &rt.ReportRetriever{Name: rc.CLICKHOUSE, Retriever: clickhouse.Init()}
		case rc.MONGO_LIVE:
			rr = &rt.ReportRetriever{Name: rc.MONGO_LIVE, Retriever: mongolive.Init()}
		case rc.MONGO_SNAP:
			rr = &rt.ReportRetriever{Name: rc.MONGO_SNAP, Retriever: mongosnap.Init()}
		case rc.REDIS:
			rr = &rt.ReportRetriever{Name: rc.REDIS, Retriever: redis.Init()}
		}
	}

	return &rt.ReportService{
		MongClient:      hr.MongClient,
		DefaultMongoDb:  hr.DefaultMongoDb,
		Logr:            hr.Logr,
		Clickhouse:      hr.ClickhouseClient,
		ReportRetriever: rr,
	}
}

func GetReport(rs *rt.ReportService, reportArgs *rt.GetActivityReportArgs) (*types.ReportApiResponse, error) {
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

	//fetching retriever
	ret := rs.ReportRetriever.Retriever

	//get data
	activityData, err := ret.GetData(rs, reportArgs)
	if err != nil {
		response.Message = err.Error()
		response.HttpStatus = http.StatusInternalServerError
		return &response, err
	}

	//format data and total
	records, totals := GetFormattedReportResponse(activityData, reportArgs.Metrics)

	msg := fmt.Sprintf("Reports fetched successfully from %s to %s", reportArgs.Start.Format("2006-01-02"), reportArgs.End.Format("2006-01-02"))
	rs.Logr.Info(context.Background(), msg)

	//set response
	response.Message = msg
	response.Success = true
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

func validate(start, end time.Time, groupBy []string, metrics []string, rs *rt.ReportService) error {
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

	if rs.ReportRetriever == nil {
		return errors.New("please provide a valid [ \"source\" ] parameter.")
	}

	return nil
}
