package acitvity

import "github.com/fuzailAhmad123/test_report/module/model"

type CreateActivityArgs struct {
	Bets  float64 `json:"bets" validate:"omitempty,min=0"`
	Wins  float64 `json:"wins" validate:"omitempty,min=0"`
	ADID  string  `json:"ad_id" validate:"required,mongoid"`
	OrgID string  `json:"org_id" validate:"required,mongoid"`
	Date  string  `json:"date" validate:"required"`
}

type CreateActivityServiceResponse struct {
	Message string         `json:"message,omitempty"`
	Data    model.Activity `json:"data,omitempty"`
}
