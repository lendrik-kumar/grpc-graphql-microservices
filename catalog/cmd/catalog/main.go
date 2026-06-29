package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/lendrik-kumar/graphql-grpc-go-microservices/catalog"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	var r catalog.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		r, err = catalog.NewElasticRepository(cfg.DatabaseURL)
		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
			return err
		}
		return nil	
	})
	defer r.Close()
	log.Println("Connected to database")
	s := catalog.NewService(r)
	log.Println("Starting gRPC server on port :8080")
	err = catalog.ListenAndServeGRPC(s, ":8080")
	if err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}