package acitvity

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fuzailAhmad123/test_report/infra/mongodb"
	"github.com/fuzailAhmad123/test_report/lib"
	"github.com/fuzailAhmad123/test_report/module/model"
	"github.com/fuzailAhmad123/test_report/module/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateActivityService a service used to create an activity.
func CreateActivityService(props *types.HTTPAPIResource, args *CreateActivityArgs) (*types.ApiResponse, error) {
	response := types.ApiResponse{
		Success:    false,
		Data:       nil,
		HttpStatus: http.StatusBadRequest,
	}

	//validations
	if args.Bets == 0 && args.Wins == 0 {
		response.Message = "Both bets and wins can't be empty"
		return &response, fmt.Errorf("Both bets and wins can't be empty")
	}

	actDate, err := lib.GetParsedTime(args.Date)
	if err != nil {
		response.Message = err.Error()
		return &response, err
	}

	//future dates not allowed
	if actDate.After(time.Now()) {
		response.Message = "Date cannot be in the future"
		return &response, fmt.Errorf("date cannot be in the future")
	}

	actID := primitive.NewObjectID()
	act := &model.Activity{
		ID:    actID,
		Bets:  args.Bets,
		Wins:  args.Wins,
		AdID:  mongodb.GetOptimisticObjectIdFromHex(args.ADID),
		OrgID: mongodb.GetOptimisticObjectIdFromHex(args.OrgID),
		Date:  actDate,
	}

	//insert into mogodb.
	_, actErr := model.InsertOne[model.Activity](context.Background(), props.DefaultMongoDb, nil, act)
	if actErr != nil {
		response.Message = actErr.Error()
		response.HttpStatus = http.StatusInternalServerError
		return &response, fmt.Errorf("Internal Server Error: ", actErr)
	}

	//insert into clickhouse.
	actClkError := InsertIntoClickhouse(props, args, actID.Hex())
	if actClkError != nil {
		response.Message = actClkError.Error()
		response.HttpStatus = http.StatusInternalServerError
		return &response, fmt.Errorf("Internal Server Error: ", actClkError)
	}

	//insert into clickhouse.
	actRedisError := InsertActivityInRedis(props, args, actDate, 1)
	if actRedisError != nil {
		response.Message = actClkError.Error()
		response.HttpStatus = http.StatusInternalServerError
		return &response, fmt.Errorf("Internal Server Error: ", actClkError)
	}

	mssg := fmt.Sprintf("Activity created with (ID:%s) successfully!", actID.Hex())
	props.Logr.Info(context.Background(), mssg)
	response.Message = fmt.Sprintf(mssg)
	response.Data = *act
	return &response, nil
}

func InsertIntoClickhouse(props *types.HTTPAPIResource, args *CreateActivityArgs, actID string) error {
	tx, err := props.ClickhouseClient.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %s", err.Error())
	}

	// Prepare the insert statement
	query := "INSERT INTO report_practice.activities (_id, org_id, ad_id, bets, wins, date) VALUES (?, ?, ?, ?, ?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %s", err.Error())
	}

	// Execute the insert
	_, err = stmt.Exec(actID, args.OrgID, args.ADID, args.Bets, args.Wins, args.Date)
	if err != nil {
		return fmt.Errorf("Insert failed: %s", err.Error())
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("Failed to commit transaction: %s", err.Error())
	}

	return nil
}

func InsertActivityInRedis(props *types.HTTPAPIResource, args *CreateActivityArgs, actDate time.Time, sign float64) error {
	//Key that will group_by all activities based on org_id and dates.
	key := fmt.Sprintf("test_actr:%s:%s", args.OrgID, actDate.Format("2006-01-02")) //test_activity_report

	fieldMapping := bson.M{
		"bets": "b",
		"wins": "w",
	}

	pipe := props.RedisClient.Client.Pipeline()

	activitySuffix := fmt.Sprintf("%s:%s", args.ADID, actDate.Format("2006-01-02"))

	if args.Bets != 0 {
		field := fmt.Sprintf("%s:%s", fieldMapping["bets"], activitySuffix)
		//update the value of "bets" with positive/negative sign
		pipe.Pipeline().HIncrByFloat(context.Background(), key, field, args.Bets*sign)
	}

	if args.Wins != 0 {
		field := fmt.Sprintf("%s:%s", fieldMapping["wins"], activitySuffix)
		//update the value of "wins" with positive/negative sign
		pipe.Pipeline().HIncrByFloat(context.Background(), key, field, args.Wins*sign)
	}

	// exec the pipeline
	_, err := pipe.Exec(context.Background())
	if err != nil {
		return fmt.Errorf("Error occured while storing activity in redis: %s", err.Error())
	}

	return nil
}
