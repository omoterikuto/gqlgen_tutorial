package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gqlgen_dataloader/datamodel"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"gqlgen_dataloader/graph/dataloader"
	"gqlgen_dataloader/graph/generated"
	"gqlgen_dataloader/graph/resolver"
	"gqlgen_dataloader/graph/storage"
)

const defaultPort = "8080"

func main() {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Logger = db.Logger.LogMode(logger.Info)

	err = db.AutoMigrate(&datamodel.User{}, &datamodel.Todo{})
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	// instantiate the DB client
	db := storage.NewMemoryStorage()
	// make a data loader
	loader := dataloader.NewDataLoader(db)
	// instantiate the gqlgen Graph Resolver
	graphResolver := resolver.NewResolver(db)
	// create the query handler
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graphResolver}))
	// wrap the query handler with middleware to inject dataloader
	dataloaderSrv := dataloader.Middleware(loader, srv)
	// register the query endpoint
	http.Handle("/query", dataloaderSrv)
	// register the playground
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	// boot the server
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
