package redis

import (

	//report constants
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

func (rt *RedisRetriever) GetData(rs *rt.ReportService, reportArgs *rt.GetActivityReportArgs) ([]model.ActivityReport, error) {
	// selectedMetrics := lo.Filter(reportArgs.Metrics, func(x string, index int) bool {
	// 	return lo.Contains(rc.ALLOWED_METRICS_FOR_REPORTING, x)
	// })

	// if len(selectedMetrics) == 0 {
	// 	fmt.Println("No valid metrics selected")
	// 	return nil, errors.New("no valid metrics selected")
	// }

	// // groupByFields := lo.Filter(reportArgs.GroupBy, func(x string, index int) bool {
	// // 	return lo.Contains(rc.ALLOWED_GROUP_BY_FOR_REPORTING, x)
	// // })

	// orgDate := reportArgs.Start.Format("2006-01-02") // Assuming daily data
	// redisKey := fmt.Sprintf("actr:%s:%s", reportArgs.OrgID, orgDate)

	// data, err := rs.GetActivityRedisData(redisKey)
	// if err != nil {
	// 	fmt.Println("Error fetching data from Redis:", err)
	// 	return nil, err
	// }

	// // Processing and grouping results
	var activities []model.ActivityReport

	// for field, value := range data {
	// 	keyParts := strings.Split(field, ":")
	// 	if len(keyParts) < 6 {
	// 		continue // Skip malformed data
	// 	}

	// 	// Extract key fields
	// 	id := keyParts[0]
	// 	orgid := keyParts[1]
	// 	adid := keyParts[2]
	// 	date := keyParts[5]

	// 	bets, _ := strconv.ParseFloat(keyParts[3], 64)
	// 	wins, _ := strconv.ParseFloat(keyParts[4], 64)
	// 	// Create activity object
	// 	obj := model.ActivityReport{
	// 		ID:    mongodb.GetOptimisticObjectIdFromHex(id),
	// 		OrgID: mongodb.GetOptimisticObjectIdFromHex(orgid),
	// 		AdID:  mongodb.GetOptimisticObjectIdFromHex(adid),
	// 		Bets:  bets,
	// 		Wins:  wins,
	// 		Date:  date,
	// 	}

	// 	activities = append(activities, obj)
	// }

	return activities, nil
}
