package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_URL"`
	CatalogURL string `envconfig:"CATALOG_URL"`
	OrderURL   string `envconfig:"ORDER_URL"`
}

func main() {
	var config AppConfig

	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatalf("Failed to process env config: %v", err)
	}

	s, err := NewGraphQLServer(config.AccountURL, config.CatalogURL, config.OrderURL)
	if err != nil {
		log.Fatalf("Failed to create GraphQL server: %v", err)
	}

	http.Handle("/graphql", handler.New(s.ToExecutableSchema()))
	http.Handle("/play", playground.Handler("sanket", "/graphql"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}