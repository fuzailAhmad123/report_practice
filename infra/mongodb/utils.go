package mongodb

import (
	"fmt"

	rc "github.com/fuzailAhmad123/test_report/module/constants" //report constants
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Only when you are sure that the id is valid and used for getting the object from the database.
func GetOptimisticObjectIdFromHex(id string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	return objID
}

func ConvertStringToObjectIds(Ids []string) []primitive.ObjectID {
	objIds := []primitive.ObjectID{}
	objIds = lo.Map(Ids, func(x string, index int) primitive.ObjectID {
		objID, err := primitive.ObjectIDFromHex(x)
		if err != nil {
			panic(err)
		}
		return objID
	})
	return objIds
}

func MakeGroupBy(groupBy []string) bson.D {
	groupBy = lo.Uniq(groupBy)
	group := bson.D{}
	for _, field := range groupBy {
		switch field {
		case rc.DATE:
			group = append(group, bson.E{Key: field, Value: bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$date"}}})
		default:
			group = append(group, bson.E{Key: field, Value: fmt.Sprintf("$%s", field)})
		}
	}

	return bson.D{{Key: "_id", Value: group}}
}
