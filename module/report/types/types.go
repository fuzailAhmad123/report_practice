package report_types

import (
	"context"
	"database/sql"
	"time"

	"github.com/fuzailAhmad123/test_report/infra/mongodb"
	"github.com/fuzailAhmad123/test_report/infra/redis"
	"github.com/fuzailAhmad123/test_report/module/model"
	"github.com/trackier/igaming-go-utils/lib/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportRetriever struct {
	Name      string
	Retriever RetrieverI
}

type RetrieverI interface {
	// functions that structs of this interface will have.
	GetCollectionName() string

	GetData(*ReportService, *GetActivityReportArgs) ([]model.ActivityReport, error)
}

type ReportService struct {
	MongClient      *mongodb.MongoClient
	DefaultMongoDb  *mongodb.MongoDefaultDatabase
	Context         context.Context
	ReportRetriever *ReportRetriever
	Clickhouse      *sql.DB
	Redis           *redis.RedisClient
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

func (rs *ReportService) GetActivityRedisData(key string) (map[string]string, error) {
	ctx := context.Background()
	data, err := rs.Redis.Client.HGetAll(ctx, key).Result()
	return data, err
}
