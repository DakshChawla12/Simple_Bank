package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/DakshChawla/simplebank/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("no authorization header"))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("no authorization header"))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationTypeBearer != authorizationType {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("authorization header type not allowed"))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
