package clickhouse

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/ClickHouse/clickhouse-go"
)

type ClickHouseClient struct {
	client *sql.DB
}

func NewClickhouseClient() (*sql.DB, error) {
	url := os.Getenv("CLICKHOUSE_URL")

	client, err := sql.Open("clickhouse", url)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
		client.Close()
		return nil, err
	}

	fmt.Println("Connected to Clickhouse")
	return client, nil
}

// Function to close clickhouse client.
func (c *ClickHouseClient) Close() {
	err := c.client.Close()
	if err != nil {
		fmt.Println("Failed to disconnect from ClickHouse:", err)
		return
	}
	fmt.Println("Disconnected from ClickHouse")
}
