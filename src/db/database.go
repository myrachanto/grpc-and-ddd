package db

import (
	"context"
	"log"

	// "go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "github.com/myrachanto/asokomonolith/support"
)

var (
	IndexRepo mongorepo
	ctx       = context.TODO()
	Mongodb   *mongo.Database
	CustomDb  *mongo.Database
	Status    = false
)

type Db struct {
	Mongohost   string `mapstructure:"Mongohost"`
	MongodbName string `mapstructure:"MongodbName"`
}
type mongorepo struct {
}

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") //locally hosted db accessed by no dockerised app
	// if clientOptions.AppName == nil {
	// 	clientOptions = options.Client().ApplyURI("mongodb://host.docker.internal:27017") //locally hosted db accessed by dockerized app
	// }
	// clientOptions := options.Client().ApplyURI("mongodb://host.docker.internal:27017")//locally hosted db accessed by dockerized app
	// clientOptions := options.Client().ApplyURI("mongodb://mongo-database:27017") //dockerized docker compose  db accessed by dockerized app
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	Mongodb = client.Database("grpcgateway")
}
