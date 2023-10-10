package main

import (
	"log"

	"github.com/myrachanto/grpcgateway/src/routes"
)

func init() {
	log.SetPrefix("gRPC: ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func main() {
	log.Println("Server started")
	routes.ApiLoader()
}
