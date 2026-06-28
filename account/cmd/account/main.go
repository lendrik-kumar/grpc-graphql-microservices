package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/lendrik-kumar/graphql-grpc-go-microservices/account"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		r, err = account.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Printf("failed to connect to database: %v", err)
		}
		return err
	})
	defer r.Close()
	log.Println("connected to database")
	s := account.NewService(r)
	log.Fatal(account.ListenAndServeGRPC(s, 8080))
}
