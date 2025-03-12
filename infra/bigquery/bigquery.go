package bigquery

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/option"
)

type BigQueryClient struct {
	Client *bigquery.Client
}

// NewBigQueryClient initializes and returns a new BigQuery client
func NewBigQueryClient() (*BigQueryClient, error) {
	ctx := context.Background()

	// Read project ID and credentials file from env variables
	projectID := os.Getenv("BIGQUERY_PROJECT_ID")
	credentialsFile := os.Getenv("BIGQUERY_CREDENTIALS_FILE")

	if projectID == "" {
		return nil, fmt.Errorf("BIGQUERY_PROJECT_ID not found")
	}
	if credentialsFile == "" {
		return nil, fmt.Errorf("BIGQUERY_CREDENTIALS_FILE not found")
	}

	client, err := bigquery.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, fmt.Errorf("failed to create BigQuery client: %v", err)
	}

	fmt.Println("Connected to BigQuery")
	return &BigQueryClient{Client: client}, nil
}

// Close closes the BigQuery client
func (bq *BigQueryClient) Close() {
	if err := bq.Client.Close(); err != nil {
		log.Println("Failed to close BigQuery client:", err)
		return
	}
	fmt.Println("Disconnected from BigQuery")
}
