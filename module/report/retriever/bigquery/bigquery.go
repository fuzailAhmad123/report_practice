package bigquery

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cloud.google.com/go/bigquery"
	rc "github.com/fuzailAhmad123/test_report/module/constants" //report constants
	"github.com/fuzailAhmad123/test_report/module/model"
	rt "github.com/fuzailAhmad123/test_report/module/report/types" //report types
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BigQueryRetriever struct{}

func Init() rt.RetrieverI {
	rt := BigQueryRetriever{}

	return &rt
}

func (rt *BigQueryRetriever) GetCollectionName() string {
	return "big-query-453307.report_practice.activities"
}

func (rt *BigQueryRetriever) GetData(rs *rt.ReportService, reportArgs *rt.GetActivityReportArgs) ([]model.ActivityReport, error) {
	selectedMetrics := lo.Filter(reportArgs.Metrics, func(x string, _ int) bool {
		return lo.Contains(rc.ALLOWED_METRICS_FOR_REPORTING, x)
	})

	if len(selectedMetrics) == 0 {
		return nil, errors.New("no valid metrics selected")
	}

	groupByFields := lo.Filter(reportArgs.GroupBy, func(x string, _ int) bool {
		return lo.Contains(rc.ALLOWED_GROUP_BY_FOR_REPORTING, x)
	})

	query := "SELECT "

	for i, field := range groupByFields {
		if field == rc.DATE {
			groupByFields[i] = "DATE(date)"
		}
	}

	if len(groupByFields) > 0 {
		query += strings.Join(groupByFields, ", ") + ", "
	}

	metricsAgg := []string{}
	for _, metric := range selectedMetrics {
		metricsAgg = append(metricsAgg, fmt.Sprintf("SUM(%s) AS %s", metric, metric))
	}

	query += strings.Join(metricsAgg, ", ")
	query += fmt.Sprintf(" FROM %s WHERE org_id = @org_id AND date BETWEEN @start AND @end", rt.GetCollectionName())

	// Prepare query parameters
	params := []bigquery.QueryParameter{
		{Name: rc.ORG_ID, Value: reportArgs.OrgID},
		{Name: rc.START, Value: reportArgs.Start.Format("2006-01-02")},
		{Name: rc.END, Value: "2025-03-12"},
	}

	if len(reportArgs.CampaignIds) > 0 {
		query += " AND ad_id IN UNNEST(@campaign_ids)"
		campaignIDs := lo.Map(reportArgs.CampaignIds, func(id primitive.ObjectID, _ int) string {
			return id.Hex()
		})
		params = append(params, bigquery.QueryParameter{Name: "campaign_ids", Value: campaignIDs})
	}

	if len(groupByFields) > 0 {
		query += " GROUP BY " + strings.Join(groupByFields, ", ")
	}

	// Run the query
	fmt.Println("Executing BigQuery:", query)

	q := rs.BigQuery.Client.Query(query)
	q.Parameters = params

	it, err := q.Read(context.Background())
	if err != nil {
		fmt.Println("BigQuery Execution Error:", err)
		return nil, err
	}

	activities, err := model.ConvertToBigQueryActivityJSON(it)
	if err != nil {
		fmt.Println("ClickHouse data converting to JSON error:", err)
		return nil, err
	}

	fmt.Println("[GetActivityReport] Data successfully fetched from BigQuery...")
	return activities, nil
}
