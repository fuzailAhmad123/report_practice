package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/fuzailAhmad123/test_report/infra/mongodb"
	rc "github.com/fuzailAhmad123/test_report/module/constants" //report constants
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
				a.Date = val.(time.Time).Format("2025-03-01")
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
