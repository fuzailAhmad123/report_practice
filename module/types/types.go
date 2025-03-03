package types

import (
	"github.com/fuzailAhmad123/test_report/infra/mongodb"
	"github.com/trackier/igaming-go-utils/lib/logger"
)

// GENERAL: Api request and response
type HTTPAPIResource struct {
	DefaultMongoDb *mongodb.MongoDefaultDatabase
	MongClient     *mongodb.MongoClient
	Logr           *logger.CustomLogger
	OrgId          string
}

type ApiResponse struct {
	Success    bool   `json:"success" validate:"required"`
	Message    string `json:"message,omitempty"`
	HttpStatus int    `json:"status,omitempty"`
	Data       any    `json:"data,omitempty"`
}

// Report arguments
type ReportRequestArgs struct {
	OrgId   string `json:"org_id" schema:"org_id" validate:"required,mongoid"`
	AdId    string `json:"ad_id" schema:"ad_id"`
	Start   string `json:"start" schema:"start" validate:"required"`
	End     string `json:"end" schema:"end" validate:"required"`
	GroupBy string `json:"groupby" schema:"groupby"`
	Metrics string `json:"metrics" schema:"metrics"`
}

type ReportApiResponse struct {
	Success    bool           `json:"success" validate:"required"`
	Message    string         `json:"message,omitempty"`
	HttpStatus int            `json:"status,omitempty"`
	Data       ReportResponse `json:"data,omitempty"`
}

type ReportResponse struct {
	Start   string   `json:"start"`
	End     string   `json:"end"`
	GroupBy []string `json:"groupby"`
	Metrics []string `json:"metrics"`
	Report  Report   `json:"report"`
}

type Report struct {
	Columns []string      `json:"columns,omitempty"`
	Records []RowFormat   `json:"records,omitempty"`
	Total   []TotalFormat `json:"total,omitempty"`
}

type RowFormat struct {
	ID    string  `json:"_id"`
	Bets  float64 `json:"bets,omitempty"`
	Wins  float64 `json:"wins,omitempty"`
	AdID  string  `json:"ad_id"`
	OrgID string  `json:"org_id"`
	Date  string  `json:"date"`
}

type TotalFormat struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}

// aggregate activity snap
type ActivitySnapAggregateArgs struct {
	OrgId   string `json:"org_id" validate:"required,mongoid"`
	Date    string `json:"date" validate:"required"`
	Refresh bool   `json:"refresh" validate:"required"`
}
