package middle

import (
	"context"
	"fmt"
	"strings"

	"github.com/myrachanto/grpcgateway/src/pasetos"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizaiontype    = "bearer"
)

func GRPCAuthorization(ctx context.Context) (*pasetos.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}
	values := md.Get(authorizationHeader)
	if len(values) < 2 {
		return nil, fmt.Errorf("missing authorizaion header")
	}
	authhead := values[0]
	fields := strings.Fields(authhead)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorizaion format")
	}
	authtype := strings.ToLower(fields[0])
	if authorizaiontype != authtype {
		return nil, fmt.Errorf("unsupported authorizaion type %s ", authtype)
	}
	accessToken := fields[1]
	paseto, err := pasetos.NewPasetoMaker()
	if err != nil {
		return nil, fmt.Errorf("authorizaion error:  %s ", err)
	}
	payload, err := paseto.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to verify token :  %s ", err)
	}
	return payload, nil
}
