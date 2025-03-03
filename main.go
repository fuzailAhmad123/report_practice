package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fuzailAhmad123/test_report/infra/mongodb"
	"github.com/fuzailAhmad123/test_report/module"
	"github.com/fuzailAhmad123/test_report/module/acitvity"
	"github.com/fuzailAhmad123/test_report/module/types"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/trackier/igaming-go-utils/lib/logger"
)

func main() {
	startServer()
}

func startServer() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error occured while loading env variables.")
	}

	logFolder := os.Getenv("REPORT_LOG_PATH")
	loggerConfig := logger.LoggerConfig{
		TimeFormat: "02-01-2006 15:04:05",
		SinkType:   logger.FILE,
		FileSinkConfig: &logger.LoggerFileSinkConfig{
			FilePath:   fmt.Sprintf("%s/logs/app.log", logFolder),
			MaxSize:    1,
			MaxBackups: 2,
			MaxAge:     1,
			Compress:   true,
		},
		BatchSize:    10,
		FlushTimeout: 5 * time.Second,
	}

	logr, err := logger.NewCustomLogger(loggerConfig)
	if err != nil {
		log.Fatalf("Failed to create custom logger: %s", err.Error())
	}

	//connection to mongodb.
	mongoClient, mongoConnectionErr := mongodb.ConnectWithMongoDb(os.Getenv("MONGODB_URL"))
	if mongoConnectionErr != nil {
		fmt.Println(mongoConnectionErr)
		os.Exit(1)
	}

	hr := &types.HTTPAPIResource{
		DefaultMongoDb: mongoClient.NewDatabase(os.Getenv("DB_NAME")),
		MongClient:     mongoClient,
		Logr:           logr,
	}

	//routes
	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK!"))
	})

	//create activity route
	r.Post("/activity", acitvity.CreateActivityController(hr))

	r.Mount("/report", module.Route(hr))

	logr.Info(context.Background(), "Server is listening on port "+os.Getenv("SERVER_PORT"))
	fmt.Println("Server is listening on port " + os.Getenv("SERVER_PORT"))
	http.ListenAndServe(os.Getenv("SERVER_PORT"), r)

	// Graceful shutdown
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		for sig := range ch {
			fmt.Printf("[INFO] %v Signal was received. Closing connected infra.\n", sig)
			mongoClient.Close()
			os.Exit(0)
		}
	}()

}
