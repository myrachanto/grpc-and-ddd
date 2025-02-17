package middle

import (
	"context"
	"fmt"
	"strings"

	"github.com/myrachanto/grpcgateway/src/pasetos"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationtype   = "bearer"
)

var log = logrus.New()

func GRPCAuthorization(ctx context.Context) (*pasetos.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.WithFields(logrus.Fields{
			"error": "missing metadata",
		}).Error("Authorization failed")
		return nil, fmt.Errorf("missing metadata")
	}
	values := md.Get(authorizationHeader)
	if len(values) < 2 {
		log.WithFields(logrus.Fields{
			"error": "missing authorization header",
		}).Error("Authorization failed")
		return nil, fmt.Errorf("missing authorization header")
	}
	authhead := values[0]
	fields := strings.Fields(authhead)
	if len(fields) < 2 {
		log.WithFields(logrus.Fields{
			"error": "invalid authorization header format",
		}).Error("Authorization failed")
		return nil, fmt.Errorf("invalid authorization header format")
	}
	authType := strings.ToLower(fields[0])
	if authorizationtype != authType {
		log.WithFields(logrus.Fields{
			"error":       "unsupported authorization type",
			"actual_type": authType,
		}).Error("Authorization failed")
		return nil, fmt.Errorf("unsupported authorization type %s ", authType)
	}
	accessToken := fields[1]
	paseto, err := pasetos.NewPasetoMaker()
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Authorization failed: could not create Paseto maker")
		return nil, fmt.Errorf("could not create Paseto maker:  %s ", err)
	}
	payload, err := paseto.VerifyToken(accessToken)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Authorization failed: token verification failed")
		return nil, fmt.Errorf("token verification failed :  %s ", err)
	}
	log.WithFields(logrus.Fields{
		"username": payload.Username,
	}).Info("Authorization successful")
	return payload, nil
}
