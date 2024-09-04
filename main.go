package main

import (
	"log"

	"github.com/myrachanto/grpcgateway/src/db"
	"github.com/myrachanto/grpcgateway/src/routes"
)

func init() {
	log.SetPrefix("gRPC: ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// @title GRPC Gateway API Documention
// @version 1.0
// @description This is a GRPC Gateway API Documention server.

// @contact.name API Support
// @contact.url https://www.chantosweb.com
// @contact.email myrachanto1@gmail.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	log.Println("Server started")
	mongodb, err := db.DbConnection()
	if err != nil {
		log.Fatal(err)
	}
	routes.ApiLoader(mongodb)
}
