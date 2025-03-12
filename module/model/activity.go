package model

import (
	"database/sql"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/fuzailAhmad123/test_report/infra/mongodb"
	rc "github.com/fuzailAhmad123/test_report/module/constants" //report constants
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/iterator"
)

func (a Activity) TableName() string {
	return "activities"
}

func (a *ActivityReport) ConvertBSONToModel(bsonData interface{}) error {
	dataBytes, err := bson.Marshal(bsonData)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(dataBytes, a)
	if err != nil {
		return err
	}
	return nil
}

type Activity struct {
	ID    primitive.ObjectID `bson:"_id" json:"_id"`
	Bets  float64            `bson:"bets" json:"bets,omitempty"`
	Wins  float64            `bson:"wins" json:"wins,omitempty"`
	AdID  primitive.ObjectID `bson:"ad_id" json:"ad_id"`
	OrgID primitive.ObjectID `bson:"org_id" json:"org_id"`
	Date  time.Time          `bson:"date" json:"date"`
}

type ActivityReport struct {
	ID    primitive.ObjectID `bson:"_id" json:"_id"`
	Bets  float64            `bson:"bets" json:"bets,omitempty"`
	Wins  float64            `bson:"wins" json:"wins,omitempty"`
	AdID  primitive.ObjectID `bson:"ad_id" json:"ad_id"`
	OrgID primitive.ObjectID `bson:"org_id" json:"org_id"`
	Date  string             `bson:"date" json:"date"`
}

func ConvertToClickhouseActivityJSON(rows *sql.Rows) ([]ActivityReport, error) {
	var activities []ActivityReport

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return activities, err
	}

	// Create a slice of interface{} to hold the scanned values
	values := make([]interface{}, len(columns))
	for i := range values {
		var v interface{}
		values[i] = &v
	}

	for rows.Next() {
		// Scan the values into the interface{} slice
		if err := rows.Scan(values...); err != nil {
			return activities, err
		}
		defer func() ([]ActivityReport, error) {
			if r := recover(); r != nil {
				fmt.Println("ConvertToClickhouseActivityJSON:", r)
				return activities, fmt.Errorf("error while converting ClickhouseActivity to JSON")
			}
			return activities, nil
		}()

		// Create a new ClickhouseActivity and populate its fields
		a := ActivityReport{}
		for i, col := range columns {
			// Use type assertion to extract the value from the interface{}
			val := *(values[i].(*interface{}))
			if val == nil {
				continue
			}
			switch col {
			case "_id":
				a.ID = mongodb.GetOptimisticObjectIdFromHex(val.(string))
			case "org_id":
				a.OrgID = mongodb.GetOptimisticObjectIdFromHex(val.(string))
			case "ad_id":
				a.AdID = mongodb.GetOptimisticObjectIdFromHex(val.(string))
			case "bets":
				a.Bets = val.(float64)
			case "wins":
				a.Wins = val.(float64)
			case "date":
				switch v := val.(type) {
				case string:
					parsedTime, err := time.Parse("2006-01-02 15:04:05", v)
					if err != nil {
						fmt.Println("Error parsing date:", err)
						continue
					}
					a.Date = parsedTime.Format("2006-01-02")

				case time.Time:
					a.Date = v.Format("2006-01-02")

				default:
					fmt.Println("Unexpected date format:", v)
				}
			}
		}

		// Append the populated activity to the activities slice
		activities = append(activities, a)
	}

	if err := rows.Err(); err != nil {
		return activities, err
	}

	return activities, nil
}

func ConvertToBigQueryActivityJSON(iter *bigquery.RowIterator) ([]ActivityReport, error) {
	var activities []ActivityReport

	for {
		var row map[string]bigquery.Value
		err := iter.Next(&row)

		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("BigQuery iteration error: %v", err)
		}

		a := ActivityReport{}

		if val, ok := row["_id"].(string); ok && val != "" {
			a.ID = mongodb.GetOptimisticObjectIdFromHex(val)
		}

		if val, ok := row["org_id"].(string); ok && val != "" {
			a.OrgID = mongodb.GetOptimisticObjectIdFromHex(val)
		}

		if val, ok := row["ad_id"].(string); ok && val != "" {
			a.AdID = mongodb.GetOptimisticObjectIdFromHex(val)
		}

		if val, ok := row["bets"]; ok {
			switch v := val.(type) {
			case int64:
				a.Bets = float64(v)
			case float64:
				a.Bets = v
			default:
				fmt.Println("Unexpected type for bets:", v)
			}
		}

		if val, ok := row["wins"]; ok {
			switch v := val.(type) {
			case int64:
				a.Wins = float64(v)
			case float64:
				a.Wins = v
			default:
				fmt.Println("Unexpected type for wins:", v)
			}
		}

		if val, ok := row["f0_"]; ok {
			switch v := val.(type) {
			case time.Time:
				a.Date = v.Format("2006-01-02")
			case string:
				parsedTime, err := time.Parse("2006-01-02", v)
				if err != nil {
					fmt.Println("Error parsing date:", err)
					continue
				}
				a.Date = parsedTime.Format("2006-01-02")
			case civil.Date:
				a.Date = fmt.Sprintf("%04d-%02d-%02d", v.Year, v.Month, v.Day)
			default:
				fmt.Println("Unexpected date format:", v)
			}
		}

		activities = append(activities, a)
	}

	return activities, nil
}

func (ar *ActivityReport) GetField(key string) string {
	switch key {
	case rc.AD_ID:
		return ar.AdID.Hex()
	case rc.ORG_ID:
		return ar.OrgID.Hex()
	case rc.DATE:
		return ar.Date
	default:
		return ""
	}
}
