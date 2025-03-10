package clickhouse

import (
	"errors"
	"fmt"
	"strings"

	rc "github.com/fuzailAhmad123/test_report/module/constants" //report constants
	"github.com/fuzailAhmad123/test_report/module/model"
	rt "github.com/fuzailAhmad123/test_report/module/report/types" //report types
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClickhouseRetriever struct{}

func Init() rt.RetrieverI {
	rt := ClickhouseRetriever{}

	return &rt
}

func (rt *ClickhouseRetriever) GetCollectionName() string {
	return "report_practice.activities"
}

func (rt *ClickhouseRetriever) GetData(rs *rt.ReportService, reportArgs *rt.GetActivityReportArgs) ([]model.ActivityReport, error) {
	selectedMetrics := lo.Filter(reportArgs.Metrics, func(x string, index int) bool {
		return lo.Contains(rc.ALLOWED_METRICS_FOR_REPORTING, x)
	})

	if len(selectedMetrics) == 0 {
		fmt.Println("No valid metrics selected")
		return nil, errors.New("no valid metrics selected")
	}

	groupByFields := lo.Filter(reportArgs.GroupBy, func(x string, index int) bool {
		return lo.Contains(rc.ALLOWED_GROUP_BY_FOR_REPORTING, x)
	})

	query := "SELECT "

	if len(groupByFields) > 0 {
		for i, field := range groupByFields {
			if field == "date" {
				groupByFields[i] = "toDate(date) AS date"
			}
		}
		query += strings.Join(groupByFields, ", ") + ", "
	}

	metricsAgg := []string{}
	for _, metric := range selectedMetrics {
		metricsAgg = append(metricsAgg, fmt.Sprintf("SUM(%s) AS %s", metric, metric))
	}

	query += strings.Join(metricsAgg, ", ")
	query += " FROM " + rt.GetCollectionName() + " WHERE org_id = ? AND date >= ? AND date <= ?"

	args := []any{
		reportArgs.OrgID,
		reportArgs.Start.Format("2006-01-02"),
		reportArgs.End.Format("2006-01-02"),
	}

	if len(reportArgs.CampaignIds) > 0 {
		query += " AND ad_id IN (?)"
		args = append(args, lo.Map(reportArgs.CampaignIds, func(id primitive.ObjectID, i int) string { return id.Hex() }))
	}

	if len(groupByFields) > 0 {
		query += " GROUP BY " + strings.Join(groupByFields, ", ")
	}

	rows, err := rs.Clickhouse.Query(query, args...)
	if err != nil {
		fmt.Println("ClickHouse Query Error:", err)
		return nil, err
	}
	defer rows.Close()

	fmt.Println("[GetActivityReport] Data successfully fetched from ClickHouse...")

	activities, err := model.ConvertToClickhouseActivityJSON(rows)
	if err != nil {
		fmt.Println("ClickHouse data converting to JSON error:", err)
		return nil, err
	}

	return activities, nil
}
