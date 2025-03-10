package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/fuzailAhmad123/test_report/infra/mongodb"
	rc "github.com/fuzailAhmad123/test_report/module/constants" //report constants
	"github.com/fuzailAhmad123/test_report/module/model"
	rt "github.com/fuzailAhmad123/test_report/module/report/types" //report types
	"github.com/samber/lo"
)

type RedisRetriever struct{}

func Init() rt.RetrieverI {
	rt := RedisRetriever{}

	return &rt
}

func (rt *RedisRetriever) GetCollectionName() string {
	return "activities"
}

func (rt *RedisRetriever) GetData(rs *rt.ReportService, reportArgs *rt.GetActivityReportArgs) ([]model.ActivityReport, error) {
	selectedMetrics := lo.Filter(reportArgs.Metrics, func(x string, index int) bool {
		return lo.Contains(rc.ALLOWED_METRICS_FOR_REPORTING, x)
	})

	if len(selectedMetrics) == 0 {
		fmt.Println("No valid metrics selected")
		return nil, errors.New("no valid metrics selected")
	}

	groupby := lo.Filter(reportArgs.GroupBy, func(x string, index int) bool {
		return lo.Contains(rc.ALLOWED_GROUP_BY_FOR_REPORTING, x)
	})

	activityRedisKey := fmt.Sprintf("test_actr:%s:%s", reportArgs.OrgID, reportArgs.Start.Format("2006-01-02"))

	data, err := rs.GetActivityRedisData(activityRedisKey)
	if err != nil {
		fmt.Print("Error while getting live activity data from redis", err.Error())
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[GetCustomerActivityFromRedis] Panic: %v", r)
			rs.Logr.Error(context.Background(), fmt.Sprintf("Error while fetching live data from Redis."))
		}
	}()

	activitiesMap := make(map[string]model.ActivityReport)
	for field, value := range data {
		keyParts := strings.Split(field, ":")
		metric := keyParts[0] // ["b" or "w"]

		//filtering based on campaigns [for org and date already filtered]
		ad_id := keyParts[1]
		if len(reportArgs.CampaignIds) > 0 && !lo.Contains(reportArgs.CampaignIds, mongodb.GetOptimisticObjectIdFromHex(ad_id)) {
			continue
		}

		date := keyParts[2]

		obj := model.ActivityReport{
			OrgID: mongodb.GetOptimisticObjectIdFromHex(reportArgs.OrgID),
			AdID:  mongodb.GetOptimisticObjectIdFromHex(ad_id),
			Date:  date,
		}

		//if bets then update bets
		if metric == "b" && lo.Contains(selectedMetrics, rc.BETS) {
			obj.Bets, _ = strconv.ParseFloat(value, 64)
		}

		//if wins then update wins
		if metric == "w" && lo.Contains(selectedMetrics, rc.WINS) {
			obj.Wins, _ = strconv.ParseFloat(value, 64)
		}

		prKey := model.GroupByKey(&obj, groupby)
		if _, ok := activitiesMap[prKey]; ok {
			//group by data
			old := activitiesMap[prKey]

			old.Bets += obj.Bets
			old.Wins += obj.Wins

			activitiesMap[prKey] = old
		} else {
			activitiesMap[prKey] = obj
		}
	}

	var activities = lo.Values(activitiesMap)
	return activities, nil
}
