package routes

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/myrachanto/grpcgateway/pb"
	"github.com/myrachanto/grpcgateway/src/api/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	docs "github.com/myrachanto/grpcgateway/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ApiLoader() {
	// go ginApiServer()
	go grpcGatewayServer()
	grpcServer()
}
func grpcGatewayServer() {
	grpcuser := users.NewUserGapiController(users.NewUserService(users.NewUserRepo()))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	grpcmux := runtime.NewServeMux()
	errd := pb.RegisterUserServiceHandlerServer(ctx, grpcmux, grpcuser)
	if errd != nil {
		log.Fatal("Error running gRPC server : ", errd)
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcmux)

	// PORT := os.Getenv("HTTP_PORT")
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal("Error running gRPC server : ", err)
	}
	log.Println("gRPC Gateway server started at :", listener.Addr().String())
	errs := http.Serve(listener, mux)
	if errs != nil {
		log.Fatal("Cannot start gRPC http gateway server : ", err)
	}
}
func grpcServer() {
	grpcserver := grpc.NewServer()
	grpcuser := users.NewUserGapiController(users.NewUserService(users.NewUserRepo()))
	pb.RegisterUserServiceServer(grpcserver, grpcuser)
	reflection.Register(grpcserver)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file in routes")
	}

	PORT := os.Getenv("GRPC_PORT")
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal("Error running gRPC server : ", err)
	}
	log.Println("gRPC server started at :", listener.Addr().String())
	errs := grpcserver.Serve(listener)
	if errs != nil {
		log.Fatal("Cannot start gRPC server : ", err)
	}
}
func ginApiServer() {

	u := users.NewUserController(users.NewUserService(users.NewUserRepo()))

	docs.SwaggerInfo.BasePath = "/"
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.Default())

	api := router.Group("/api")

	router.POST("/register", u.Create)
	router.POST("/login", u.Login)

	api.GET("/logout", u.Logout)
	api.POST("/users/shop", u.Create)
	api.GET("/users", u.GetAll)
	api.GET("/users/:code", u.GetOne)
	api.PUT("/users/password", u.PasswordUpdate)

	router.GET("/health", HealthCheck)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file in routes")
	}

	PORT := os.Getenv("HTTP_PORT")
	router.Run(PORT)
}
