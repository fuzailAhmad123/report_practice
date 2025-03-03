package clickhouse

import (
	"github.com/fuzailAhmad123/test_report/module/model"
	rt "github.com/fuzailAhmad123/test_report/module/report/types" //report types
)

type ClickhouseRetriever struct{}

func Init() rt.RetrieverI {
	rt := ClickhouseRetriever{}

	return &rt
}

func (rt *ClickhouseRetriever) GetCollectionName() string {
	return "activities"
}

func (rt *ClickhouseRetriever) GetData(rs *rt.ReportService, reportArgs *rt.GetActivityReportArgs) ([]any, error) {
	return []any{}, nil
}

func (rt *ClickhouseRetriever) ConvertToBSON(data []any) ([]model.ActivityReport, error) {
	var activities []model.ActivityReport

	//connvert clickhouse fetched data to bson
	return activities, nil
}
