package redis

import (
	"github.com/fuzailAhmad123/test_report/module/model"
	rt "github.com/fuzailAhmad123/test_report/module/report/types" //report types
)

type RedisRetriever struct{}

func Init() rt.RetrieverI {
	rt := RedisRetriever{}

	return &rt
}

func (rt *RedisRetriever) GetCollectionName() string {
	return "activities"
}

func (rt *RedisRetriever) GetData(rs *rt.ReportService, reportArgs *rt.GetActivityReportArgs) ([]any, error) {
	return []any{}, nil
}

func (rt *RedisRetriever) ConvertToBSON(data []any) ([]model.ActivityReport, error) {
	var activities []model.ActivityReport

	//connvert redis fetched data to bson
	return activities, nil
}
