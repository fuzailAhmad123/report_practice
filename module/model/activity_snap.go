package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a ActivitySnap) TableName() string {
	return "activities_snap"
}

type ActivitySnap struct {
	ID    primitive.ObjectID `bson:"_id" json:"_id"`
	Bets  float64            `bson:"bets" json:"bets,omitempty"`
	Wins  float64            `bson:"wins" json:"wins,omitempty"`
	AdID  primitive.ObjectID `bson:"ad_id" json:"ad_id"`
	OrgID primitive.ObjectID `bson:"org_id" json:"org_id"`
	Date  time.Time          `bson:"date" json:"date"`
}
