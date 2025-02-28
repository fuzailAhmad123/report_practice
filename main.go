package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/fuzailAhmad123/test_report/infra/mongodb"
	"github.com/fuzailAhmad123/test_report/module"
	"github.com/fuzailAhmad123/test_report/module/acitvity"
	"github.com/fuzailAhmad123/test_report/module/types"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	startServer()
}

func startServer() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error occured while loading env variables.")
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
	}

	//routes
	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK!"))
	})

	//create activity route
	r.Post("/activity", acitvity.CreateActivityController(hr))

	r.Mount("/report", module.Route(hr))

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
