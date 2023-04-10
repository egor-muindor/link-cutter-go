package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"cutter-url-go/internal/handlers/cut"
	"cutter-url-go/internal/handlers/redirect"
	"cutter-url-go/internal/repositories/link"
)

func main() {
	ctx := context.Background()
	mongoURI := os.Getenv("MONGODB_URI")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI))
	if err != nil {
		panic(err)
	}
	db := client.Database("ls")

	repo := link.NewRepository(db)
	ch := cut.NewHandler(repo)
	rh := redirect.NewHandler(repo)

	http.Handle("/cut", ch)

	http.Handle("/", rh)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
