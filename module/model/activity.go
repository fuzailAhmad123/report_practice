package model

import (
	"time"

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

func (c *ActivityReport) GetField(key string) string {
	switch key {
	case rc.AD_ID:
		return c.AdID.Hex()
	case rc.ORG_ID:
		return c.OrgID.Hex()
	case rc.DATE:
		return c.Date
	default:
		return ""
	}
}
