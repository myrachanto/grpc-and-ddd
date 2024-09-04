package db

import (
	"context"
	"fmt"

	// "go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "github.com/myrachanto/asokomonolith/support"
)

var (
	ctx = context.TODO()
)

func DbConnection() (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") //locally hosted db accessed by no dockerised app
	// if clientOptions.AppName == nil {
	// clientOptions := options.Client().ApplyURI("mongodb://host.docker.internal:27017") //locally hosted db accessed by dockerized app
	// }
	// clientOptions := options.Client().ApplyURI("mongodb://host.docker.internal:27017")//locally hosted db accessed by dockerized app
	// clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017") //dockerized docker compose  db accessed by dockerized app
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping the database")
	}
	Mongodb := client.Database("grpcgateway")
	return Mongodb, nil
}
