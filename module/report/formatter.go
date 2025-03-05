package report

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/fuzailAhmad123/test_report/infra/mongodb"
	"github.com/fuzailAhmad123/test_report/lib"
	rc "github.com/fuzailAhmad123/test_report/module/constants"
	"github.com/fuzailAhmad123/test_report/module/model"
	rt "github.com/fuzailAhmad123/test_report/module/report/types" //report types
	"github.com/fuzailAhmad123/test_report/module/types"
)

func ConvertReportQueryParams(args *types.ReportRequestArgs) (*rt.GetActivityReportArgs, error) {
	start, err := lib.GetParsedTime(args.Start)
	if err != nil {
		return nil, fmt.Errorf("Start date conversion error: %v", err)
	}

	end, err := lib.GetParsedTime(args.End)
	if err != nil {
		return nil, fmt.Errorf("End date conversion error: %v", err)
	}

	groupBy := lo.Map(strings.Split(args.GroupBy, ","), func(x string, index int) string {
		return strings.TrimSpace(x)
	})

	metrics := lo.Map(strings.Split(args.Metrics, ","), func(x string, index int) string {
		return strings.TrimSpace(x)
	})

	// Convert campaign IDs to MongoDB ObjectIDs
	var campaignObjectId []primitive.ObjectID
	adIds := lo.Map(strings.Split(args.AdId, ","), func(x string, index int) string { return strings.TrimSpace(x) })
	adFilteredIds := lo.Filter(adIds, func(id string, _ int) bool { return id != "" })

	if len(adFilteredIds) > 0 {
		campaignObjectId = mongodb.ConvertStringToObjectIds(adFilteredIds)
	}

	// Return parsed values in a structured format
	return &rt.GetActivityReportArgs{
		Start:       start,
		End:         end,
		GroupBy:     groupBy,
		Metrics:     metrics,
		CampaignIds: campaignObjectId,
		OrgID:       args.OrgId,
	}, nil
}

func GetFormattedReportResponse(activityData []model.ActivityReport, metrics []string) ([]types.RowFormat, []types.TotalFormat) {
	//format data into row format
	var records []types.RowFormat
	for _, ad := range activityData {
		row := types.RowFormat{
			// ID:    ad.ID.Hex(),
			Bets:  ad.Bets,
			Wins:  ad.Wins,
			OrgID: ad.OrgID.Hex(),
			AdID:  ad.AdID.Hex(),
			Date:  ad.Date,
		}
		records = append(records, row)
	}

	//format the total row
	totalMap := make(map[string]float64)
	for _, metric := range metrics {
		totalMap[metric] = 0
	}

	for _, ad := range activityData {
		for _, metric := range metrics {
			switch metric {
			case rc.BETS:
				totalMap[metric] += ad.Bets
			case rc.WINS:
				totalMap[metric] += ad.Wins
			}
		}
	}

	// creaate formatted tottal row array
	var totals []types.TotalFormat
	for _, metric := range metrics {
		totals = append(totals, types.TotalFormat{Key: metric, Value: totalMap[metric]})
	}

	return records, totals
}
