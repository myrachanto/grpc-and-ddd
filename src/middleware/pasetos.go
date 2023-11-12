package middle

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/grpcgateway/src/pasetos"
)

const (
	authorisationHeaderKey = "Authorization"
	authorisationType      = "Bearer"
)

func PasetoAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(authorisationHeaderKey)
		if len(authorizationHeader) == 0 {
			c.JSON(http.StatusUnauthorized, httperrors.NewBadRequestError("authorization header not provided"))
			return
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			c.JSON(http.StatusUnauthorized, httperrors.NewBadRequestError("Invalid Authorization format provided"))
			return
		}
		authtype := fields[0]
		if authtype != authorisationType {
			c.JSON(http.StatusUnauthorized, httperrors.NewBadRequestError("That type of Authorization is not allowed here!"))
			return
		}
		accessToken := fields[1]
		maker, _ := pasetos.NewPasetoMaker()
		payload, err := maker.VerifyToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, httperrors.NewBadRequestError("That token is invalid!"))
			return
		}
		if payload.Username == "" {
			c.JSON(http.StatusUnauthorized, httperrors.NewBadRequestError("That User is invalid!"))
			return
		}
		c.Set("username", payload.Username)
		c.Next()
	}
}
