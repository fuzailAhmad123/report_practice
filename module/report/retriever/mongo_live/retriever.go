package mongolive

import (
	"errors"
	"fmt"

	"github.com/fuzailAhmad123/test_report/infra/mongodb"
	rc "github.com/fuzailAhmad123/test_report/module/constants" //report constant
	"github.com/fuzailAhmad123/test_report/module/model"
	rt "github.com/fuzailAhmad123/test_report/module/report/types" //report types
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoLiveRetriever struct{}

func Init() rt.RetrieverI {
	rt := MongoLiveRetriever{}

	return &rt
}

func (rt *MongoLiveRetriever) GetCollectionName() string {
	return "activities"
}

func (rt *MongoLiveRetriever) GetData(rs *rt.ReportService, reportArgs *rt.GetActivityReportArgs) ([]any, error) {
	collection := rs.DefaultMongoDb.Db.Collection(rt.GetCollectionName())

	// 1. Match Stage Fields
	matchStageFields := bson.D{
		{Key: rc.ORG_ID, Value: mongodb.GetOptimisticObjectIdFromHex(reportArgs.OrgID)},
		{Key: rc.DATE, Value: bson.D{{Key: rc.GREATER_THAN_EQUALS, Value: reportArgs.Start}, {Key: rc.LESSER_THAN_EQUALS, Value: reportArgs.End}}},
	}

	if len(reportArgs.CampaignIds) > 0 {
		matchStageFields = append(matchStageFields, bson.E{Key: rc.AD_ID, Value: bson.D{{Key: rc.INCLUDES, Value: reportArgs.CampaignIds}}})
	}

	// Allowed metrics
	selectedMetrics := lo.Filter(reportArgs.Metrics, func(x string, index int) bool {
		return lo.Contains(rc.ALLOWED_METRICS_FOR_REPORTING, x)
	})

	if len(selectedMetrics) == 0 {
		fmt.Println("No valid metrics selected")
		return nil, errors.New("no valid metrics selected")
	}

	// 2. GroupBy Stage Fields
	groupby := lo.Filter(reportArgs.GroupBy, func(x string, index int) bool {
		return lo.Contains(rc.ALLOWED_GROUP_BY_FOR_REPORTING, x)
	})

	groupByStageFields := mongodb.MakeGroupBy(groupby)

	// Adding the selected metrics to the grouped data
	for _, metric := range selectedMetrics {
		groupByStageFields = append(groupByStageFields, bson.E{Key: metric, Value: bson.D{{Key: rc.SUM, Value: "$" + metric}}})
	}

	// 3. Projection Fields
	projectStageFields := bson.D{
		{Key: rc.ID, Value: 0},
	}

	// Add grouping fields to projection
	for _, field := range groupby {
		projectStageFields = append(projectStageFields, bson.E{Key: field, Value: rc.PROJECTION_PREFIX + field})
	}

	// Add selected metrics to projection
	for _, metric := range selectedMetrics {
		projectStageFields = append(projectStageFields, bson.E{Key: metric, Value: 1})
	}

	// 4. Fetch data
	result, err := mongodb.Aggregation(collection, matchStageFields, groupByStageFields, projectStageFields)
	if err != nil {
		fmt.Println("Aggregation Error:", err)
		return nil, err
	}

	fmt.Println("[GetActivityLiveData] Activties from Mongodb Live Data successfully fetched....")

	return result, nil
}

func (rt *MongoLiveRetriever) ConvertToBSON(data []any) ([]model.ActivityReport, error) {
	var activities []model.ActivityReport
	for _, obj := range data {
		activity := &model.ActivityReport{}
		err := activity.ConvertBSONToModel(obj)
		if err != nil {
			fmt.Println("Erro is: ", err)
			return nil, err
		}
		activities = append(activities, *activity)
	}

	return activities, nil
}
