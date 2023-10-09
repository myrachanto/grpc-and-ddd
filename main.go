package main

import (
	"log"

	"github.com/myrachanto/grpcgateway/src/routes"
)

func init() {
	log.SetPrefix("GRPCGateway server ...... ")
}
func main() {
	log.Println("Server started")
	routes.ApiLoader()
}
