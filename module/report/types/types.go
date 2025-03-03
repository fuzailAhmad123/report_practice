package report_types

import (
	"context"
	"time"

	"github.com/fuzailAhmad123/test_report/infra/mongodb"
	"github.com/fuzailAhmad123/test_report/module/model"
	"github.com/trackier/igaming-go-utils/lib/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportRetriever struct {
	Name      string
	Docs      []any
	Retriever RetrieverI
}

type RetrieverI interface {
	// functions that structs of this interface will have.
	GetCollectionName() string

	GetData(*ReportService, *GetActivityReportArgs) ([]any, error)

	ConvertToBSON([]any) ([]model.ActivityReport, error)
}

type ReportService struct {
	MongClient      *mongodb.MongoClient
	DefaultMongoDb  *mongodb.MongoDefaultDatabase
	Context         context.Context
	ReportRetriever *ReportRetriever
	Logr            *logger.CustomLogger
}

type GetActivityReportArgs struct {
	Start       time.Time
	End         time.Time
	GroupBy     []string
	Metrics     []string
	CampaignIds []primitive.ObjectID
	OrgID       string
}
