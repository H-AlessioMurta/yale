package main

import (
	"log"
	"net/http"
	"os"
	l "yale/borrowing/logger"
	"yale/borrowing/graph"
	"yale/borrowing/graph/generated"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	log.SetPrefix("\033[36mborrowingsvc\033[0m: ")
	log.SetOutput(os.Stdout)
	l.LogInfo("connect to http://localhost:"+port+"/ for GraphQL playground")
	http.ListenAndServe(":"+port, nil)
}
