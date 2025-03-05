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
	if actClkError != "" {
		response.Message = actClkError
		response.HttpStatus = http.StatusInternalServerError
		return &response, fmt.Errorf("Internal Server Error: ", actClkError)
	}

	mssg := fmt.Sprintf("Activity created with (ID:%s) successfully!", actID.Hex())
	props.Logr.Info(context.Background(), mssg)
	response.Message = fmt.Sprintf(mssg)
	response.Data = *act
	return &response, nil
}

func InsertIntoClickhouse(props *types.HTTPAPIResource, args *CreateActivityArgs, actID string) string {
	ctx := context.Background()
	tx, err := props.ClickhouseClient.Begin()
	if err != nil {
		props.Logr.Error(ctx, err.Error())
		return fmt.Sprintf("failed to begin transaction: %s", err.Error())
	}

	// Prepare the insert statement
	query := "INSERT INTO report_practice.activities (_id, org_id, ad_id, bets, wins, date) VALUES (?, ?, ?, ?, ?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		props.Logr.Error(ctx, err.Error())
		return fmt.Sprintf("Failed to prepare statement: %s", err.Error())
	}

	// Execute the insert
	_, err = stmt.Exec(actID, args.OrgID, args.ADID, args.Bets, args.Wins, args.Date)
	if err != nil {
		props.Logr.Error(ctx, err.Error())
		return fmt.Sprintf("Insert failed: %s", err.Error())
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		props.Logr.Error(ctx, err.Error())
		return fmt.Sprintf("Failed to commit transaction: %s", err.Error())
	}

	return ""
}
