build:
	@go build -o grpcgateway

run:
	@go run .

test:
	@go test -v ./...

testCover:
	@go test -v ./... -cover

swagger:
	@"$HOME/go/bin/swag init -g ./src/routes/routes.go"

dockerize:
	@docker build -t grpcgateway:latest .

dockerrun:
	@docker run --name grpcapp -p 4500:4500 grpcgateway:latest
protogen:
	@rm -f pb/*.go
	@protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    --grpc-gateway_opt logtostderr=true \
	--grpc-gateway_out=pb  --grpc-gateway_opt=paths=source_relative \
    proto/*.proto
evanslist:
	@evans --host localhost --port 5500 -r repl