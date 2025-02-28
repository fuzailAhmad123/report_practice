package activity_snap

import (
	"context"
	"time"

	"github.com/fuzailAhmad123/test_report/infra/mongodb"
	"github.com/fuzailAhmad123/test_report/module/model"
	"github.com/fuzailAhmad123/test_report/module/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateActivitySnapShots(rs *types.HTTPAPIResource, activities []types.RowFormat, createdDate time.Time) error {
	var snapActivities []model.ActivitySnap
	for _, act := range activities {
		actSnap := model.ActivitySnap{
			ID:    primitive.NewObjectID(),
			Bets:  act.Bets,
			Wins:  act.Wins,
			AdID:  mongodb.GetOptimisticObjectIdFromHex(act.AdID),
			OrgID: mongodb.GetOptimisticObjectIdFromHex(act.OrgID),
			Date:  createdDate,
		}

		snapActivities = append(snapActivities, actSnap)
	}

	_, err := model.InsertMany(context.Background(), rs.DefaultMongoDb, nil, snapActivities)
	if err != nil {
		return err
	}

	return nil
}
