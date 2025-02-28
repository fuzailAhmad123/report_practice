package validation

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fuzailAhmad123/test_report/module/types"
	"github.com/go-playground/validator"
	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ValidateRequestBody validates the request body
func ValidateRequestBody(r *http.Request, args any) (*types.ApiResponse, error) {
	result := types.ApiResponse{
		Success:    false,
		Data:       nil,
		HttpStatus: http.StatusBadRequest,
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&args); err != nil {
		result.Message = fmt.Sprintf("unable to decode request body: %v", err)
		return &result, err
	}

	validate := validator.New()
	if err := validate.RegisterValidation("mongoid", IsValidMongoIDStructTag); err != nil {
		result.Message = fmt.Sprintf("error registering custom validation: %v", err.Error())
		return &result, err
	}

	if err := validate.Struct(args); err != nil {
		result.Message = fmt.Sprintf("validation error: %v", err.Error())
		return &result, err
	}

	return &result, nil
}

var decoder = schema.NewDecoder()

// ValidateGetQueryParams validates the query param data.
func ValidateGetQueryParams(r *http.Request, args any) (*types.ApiResponse, error) {
	result := types.ApiResponse{
		Success:    false,
		Data:       nil,
		HttpStatus: http.StatusBadRequest,
	}

	decoder.IgnoreUnknownKeys(true)
	if err := decoder.Decode(args, r.URL.Query()); err != nil {
		result.Message = fmt.Sprintf("unable to decode query parameters: %v", err)
		return &result, err
	}

	validate := validator.New()
	if err := validate.RegisterValidation("mongoid", IsValidMongoIDStructTag); err != nil {
		result.Message = fmt.Sprintf("error registering custom validation: %v", err.Error())
		return &result, err
	}

	if err := validate.Struct(args); err != nil {
		result.Message = fmt.Sprintf("validation error: %v", err.Error())
		return &result, err
	}

	return &result, nil
}

// isValidMongoID checks if the provided ID is a valid MongoDB ObjectID .
func isValidMongoID(id string) bool {
	if len(id) == 24 {
		if _, err := primitive.ObjectIDFromHex(id); err == nil {
			return true
		}
	}
	return false
}

// IsValidMongoId validates a field to check if it is a valid MongoDB ObjectID .
func IsValidMongoIDStructTag(field validator.FieldLevel) bool {
	return isValidMongoID(field.Field().String())
}
